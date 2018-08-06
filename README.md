# LyricerSpotify
Simple console app displaying the lyrics of the currently played track on spotify.

As for now it's just proof of concept, but I'm going to make it a full fledged app in the future. 

It knows when to fetch new lyrics by how much time is remaining to the tracks end.
So if you stopped the track or just switched in in the middle type "r" en "enter" to refresh force fetching new lyrics.

# WIP
Note that it's still just a work in progress.
There is still no token refreshing implemented so you need to authorize app once again after abouat an hour.

And there are no app's client's secret and client's id so no can use it at the moment.
And because of some legal issues there is an API for lyrics fetching but no implementation for it (I have it on my side but yeah).

It's more so just showcase of my project :P

# But I want to use it anyway.

Well if you really want to.

1. Create your own app on spotify.
2. Copy `conf.json` file and save it as a `hidden_conf.json` in the packages root directory.
3. Copy yours spotidy app `client id` and `client secret` to the `hidden_conf.json`.
4. Write implementation for the `Fetcher` interface in the `lyrics` package. 
5. In the `main.go` replace `lyrics.TekstowoFetcher{}` in [line 47](https://github.com/gala377/LyricerSpotify/blob/8118232f0cce47092c4b7d7788187f9335c95aad/main.go#L47) with your own `Fetcher` implementation.

All should work now.

# Known issues.
1. `main.go` is a mess.
2. No token refreshing.
3. Fetching lyrics multiple times for the same song.
4. Fetching lyrics even if the track is stopped.
5. No console clear.
6. So much logs.
