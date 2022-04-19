# Kiki's Feed Service

This learning project (to teach myself go) is a RSS aggregator. I've built it because most feed Readers are really wonky (reading the article almost never work either because of the reader, or the website followed). Plus it contains tons of features I don't need, I just want stuff under my radar and not having to mark 400 articles on read.

This project, when ran, will fetch each rss feed you are interested and build a simple html page, listing every article polled, sorted by most recent.
<div style="text-align:center">
  <img src="https://github.com/Elesh-Norn/KikiFeedService/raw/main/examplepicture.png"/>
</div>

## Configuring it

Everything is in the yaml.config file. Rename or copy the example to config.yaml.
- adresses must be the rss feed. Must be between parenthesis.
- articlenumber: Number of article per site to display. Default is 10.
- title: The title you want your page to appear
- useragent: You can customise how you want to appear, so people recognize you.
### Installing 

Clone this repository then follow go instructions to install the code.
If you don't want to install it just run `go run .` in the folder

### Usage

`kikifeedservice` or `gun run .` without argument will just fetch article and build a static page.

`kikifeedservice -server` or `go run . -server` will launch a local server on default port (8090).

`kikifeedservice -server -port` will allows you to choose your own port

There is a lot of room to improve since it's my first time touching go. I wouldn't consider using this for anything serious in this state. I might improve it in the upcoming weeks.

I could: 
- ~~Use goroutines to fetch the site (it's basic and important go, so I need to learn it).~~ 
- ~~Make the project easily deployable online and poll once a day.~~ (Deployed this on a private 
server and ran it with a cron job. Accessible at rss.emberger.xyz).
- ~~Have a Custom User Agent.~~
