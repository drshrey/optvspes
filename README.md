Twitter Score:
-------------

Does some basic sentiment analysis on a user's tweets and spits back whether they are a "positive" or "negative" person. 

Right now, I'm using the Alchemy API which means the daily call limit is 1000 per day. I've limited the amount of tweets processed to 200, so you don't burn right through right away. 

####Here's the website to get a free key:
[Alchemy API Free API Key Page](http://www.alchemyapi.com/api/register.html)

For now, I've attached a sample json file for using with the tool. Since one would need a Twitter account for api key details and all that jazz, I've just attached this json file to use (simply to see the tool's functionality). 

Further Improvements:

Develop an independent semantic analysis library OR port NLTK's semantic libraries into Go


####All you have to do is include this file:
1. api_key.txt

This should just include your api key you received from Alchemy. 


Simple Functionality:
go run backend.go

Functionality for different user names:
coming soon

APIs/external libraries Used:

[Alchemy API Go port](https://github.com/lineback/alchemyapi_go)
[Twitter API Go port](https://github.com/kurrik/twittergo)
[Mapstructure library](https://github.com/mitchellh/mapstructure)
