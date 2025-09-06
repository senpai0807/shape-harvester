# Shape Harvester
- I wanted to upload an old Shape Harvester that I made about a year ago. It uses this package https://github.com/chromedp/chromedp to initalize and communicate with the browser. I updated it to work with browsers other than Chrome like Edge, Brave, & Opera and also updated product url to something current, I tested it a bit and it worked for about 5 requests then it stopped, so as of currently it doesn't work. I wanted to upload to give people who are trying to create their own harvester and idea of how to go about it. Potentially introducing proxy usage and a temporary directory could get this repo to work, but I'm not 100% on that

- ## Usage
```
1. git clone https://github.com/senpai0807/shape-harvester.git
2. cd shape-harvester
3. go mod tidy
 4. go run ./src/main.go
 ```

- To change browser type, simply go to the main.go and change browserType to Chrome, Edge, or Brave
 
### How Does It Work?
- Blocks Request -> Intercepts Request Header From /cart_items -> Sends To Server -> Repeats
- You will need to either remove the user of a server that the harvester sends the headers to or simply create a server and push the headers into a jar for usage
