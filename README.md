# Journey: 幫你找工作的好夥伴

> 這是一個協助求職者快速搜尋工作的網站

## Why?

厭煩於打開 104 搜尋職缺，一個一個查看職缺內容找尋適合你的工作嗎？

那你應該使用 Journey ，Journey 可以幫你摘要職缺，找出職缺的關鍵字，透過我們的網頁，你可以進階搜尋，

想要零食櫃嗎？ 打勾！

想要Code Review嗎？ 打勾！

想要彈性上下班嗎？ 打勾！

只看符合你心目中理想的職缺，不要再浪費時間，將那些沒誠意的職缺剔除掉吧！！

## How?

我們採用以下技術來達成 Journey

- Go: 一個簡單又快速的程式語言
- Colly: 靜態網頁爬蟲
- go-rod: 動態網頁爬蟲
- gorm: ORM
- postgres: 資料庫
- docker: 容器化
- godotenv: 環境變數

透過爬蟲，我們將求職網站的職缺搜集並整理，並交由深度學習模型進行職缺的關鍵字抓取，讓你可以精準求職，找到想要的工作！

## Installation

如何自行運行此服務

## Usage Example

## API

## Contributing

## TODO

- [ ] 加速每個職缺網站的內容爬取，現在的每個 thread 會負責一個頁面的爬取，並迴圈爬取該頁面底下所有職缺的內容
- [ ] 撰寫測試以及實現基本的 CI/CD
- [ ] 前端頁面實現
- [ ] 將整理的文本資料丟進 message queue 傳遞給 python 深度模型服務
- [ ] 定時爬取資料，統計技能需求歷史變化，可以用來得知哪個技能最熱門、哪個技能最貴、哪些技能同時出現
