# Alex App - AI Assistant for your Trip through the Art Gallery
I built this web app to help people explore art galleries and museums in a more interactive
way. You just name the painting and Alex will give you a quick introduction. Besides that,
you an also ask it questions to get more insights.

The basic idea is to have a digital guide that gives you a personalized tour.

## How it works under the hood?
I use OpenAI to generate the descriptions and answers to your questions. AWS Polly converts
that text to an audio file which you can play in the web app.

## Why didn't it work?
Mostly, because I suck at marketing and sales. I tried to collaborate with museums and galleries
in Germany but couldn't convince them to try it out. Also because most of them already have
apps and there is no incentive to have a competing offer.

That collaboration is crucial tho because GPT models often hallucianted especially when
asked for context for less prominent paintings. To avoid non-sensical answers, I had to
collaborate with galleries.

But what put the nail in the coffin was the fact that this app has no moat. OpenAI just started
to extend its app with audio interactions. That means, I can also just open a chat in their app
and have the same style of conversation I would have in my app.

## Project structure
* `app` - React app or web app
* `cmd` - Golang backend entrypoints
* `cdc` - Consumer Defined Contract tests
