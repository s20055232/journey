package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)


var Browser *rod.Browser

func initializeBrowser(rodURL string){
	// This example is to launch a browser remotely, not connect to a running browser remotely,
	// to connect to a running browser check the "../connect-browser" example.
	// Rod provides a docker image for beginners, run the below to start a launcher.Manager:
	//
	//     docker run --rm -p 7317:7317 ghcr.io/go-rod/rod
	//
	// For available CLI flags run: docker run --rm ghcr.io/go-rod/rod rod-manager -h
	// For more information, check the doc of launcher.Manager
	l := launcher.MustNewManaged(rodURL)

	// You can also set any flag remotely before you launch the remote browser.
	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
	l.Set("disable-gpu").Delete("disable-gpu")

	// Launch with headful mode
	l.Headless(true)

	Browser = rod.New().Client(l.MustClient()).MustConnect()
}

func getPages(rodURL, domain, jobCat string) int{
	initializeBrowser(rodURL)
	// You may want to start a server to watch the screenshots of the remote browser.
	// launcher.Open(browser.ServeMonitor(""))
	page := Browser.MustPage(fmt.Sprintf("https://%v/jobs/search/?jobcat=%v&page=1", domain, jobCat))
	page.MustWaitStable()
	pageSelectorElement, _ := page.Element(".page-select.js-paging-select")
	optionElement, _ := pageSelectorElement.Elements("option")
	pages := len(optionElement)
	return pages
}


func openJobPage(url string) *rod.Page{
	page := Browser.MustPage(url)
	page.MustWaitStable()
	return page
}

func getJobContent(page *rod.Page) string{
	jbElement, _ := page.Element(".job-address")
	jd := jbElement.MustAttribute("jobdescription")
	salary := jbElement.MustAttribute("salary")
	jobContent := fmt.Sprintf("%v\n%v", *jd, *salary)
	Logger.Println("jobContent: ", jobContent)
	return jobContent
}

func getJobDescription(url string) string{
	page := openJobPage(url)
	jobContent := getJobContent(page)
	jobRequirement := getJobRequirement(page)
	jobDescription := fmt.Sprintf("%v %v", jobContent, jobRequirement)
	return jobDescription
}

func getJobRequirement(page *rod.Page) string {
	jbElement := page.MustElement(".job-requirement")
	titles := jbElement.MustElements(".h3")
	contents := jbElement.MustElements(".t3.mb-0")
	jobRequirement := ""
	for idx, title := range titles {
        text := title.MustText()
		content := ""		
		if idx > len(contents) {
			content = jbElement.MustElement("p").MustText()
		}else{
			content = contents[idx].MustText()
		}
		jr := fmt.Sprintf("%s: %s", text, content)
        Logger.Println(jr)
		jobRequirement += fmt.Sprintf("%s\n", jr)
    }
	return jobRequirement
}

func createCollector(domain string, channel chan<- Job) *colly.Collector{
	outerCollector := colly.NewCollector(
		// colly.AllowedDomains(domain),
		colly.Async(true),
	)
	outerCollector.Limit(&colly.LimitRule{
		// DomainGlob:  "*104.com.tw*",
		Parallelism: 2,
		RandomDelay: 10 * time.Second,
	})
	outerCollector.OnRequest(
		func(r *colly.Request) { log.Println("Visiting", r.URL) })
		outerCollector.OnError(func(rp *colly.Response, err error) {
		Logger.Println("Something went wrong:", err)
		Logger.Println("Status Code:", rp.StatusCode)
		Logger.Println("Body:", string(rp.Body))
	})
	outerCollector.OnResponse(func(r *colly.Response) {
		Logger.Println("Visited", r.Request.URL)
		Logger.Println("Status Code:", r.StatusCode)
		// Logger.Println("Body:", string(r.Body))
	})
	outerCollector.OnHTML("article", func(e *colly.HTMLElement) {
		// 使用Attr方法來獲取data-job-name的屬性值
		job := Job{}
		company := Company{}
		attrs := []string{"data-job-name", "data-cust-name", "data-indcat-desc"}
		for _, attr := range attrs {
			findAttr := e.Attr(attr)
			switch attr {
			case "data-job-name":
				job.Name = findAttr
			case "data-cust-name":
				company.Name = findAttr
			case "data-indcat-desc":
				company.Industry = findAttr
			}
			Logger.Println(attr, ":", findAttr)
		}
		jobLink := e.ChildAttr("a.js-job-link", "href")
		if jobLink != "" {
			url, _ := strings.CutPrefix(jobLink, "//www.")
			jobURL := "https://" + url
			job.JobURL = jobURL
			Logger.Println("職缺連結:", jobURL)
			jd := getJobDescription(jobURL)
			job.Description = jd
		}

		companyLink := e.ChildAttr("ul.b-list-inline.b-clearfix li a", "href")
		if companyLink != "" {
			url, _ := strings.CutPrefix(jobLink, "//www.")
			companyURL := "https://" + url
			company.CompanyURL = companyURL
			Logger.Println("公司連結:", companyURL)
			// jobRequorement := 
			company.Description = ""
		}
		job.Company = company
		channel <- job
		// Logger.Panic("STOP!!")
	})
	return outerCollector
}


func scrape(domain string, pages int, channel chan<- Job, quit chan<- int) {
	collector := createCollector(domain, channel)
	// "jobcat" equals "backend" initially, but it can be changed later.
	url := fmt.Sprintf("https://%v/jobs/search/?jobcat=2007001016", domain)
	// Used for collecting jobs listed on pages.
	q, _ := queue.New(
		10, // Number of consumer threads, this means how many pages can be collected simultaneously.
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)
	for i := 2; i <= pages; i++ {
		q.AddURL(fmt.Sprintf("%v&page=%v", url, i))
	}
	q.Run(collector)
	collector.Visit(url)
	collector.Wait()
	quit <- 0
}
