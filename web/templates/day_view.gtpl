<html>
    <head>
        {{template "style.gtpl"}}
    </head>
    <body>
        <h1>{{.FileName}}</h1>
        <a href="/">Back to index</a>
        <center><div id="playerdiv"></div></center>
        <ul>
            {{range $index, $link := .Links}}
                <li><a href="#" onclick='player.loadVideoById("{{$link}}", 0, "default")'>{{index $.Titles $index}}</a>
                <a href="/delete/{{urlquery (index $.OriginalFileName $index)}}/{{$.FileName}}" onclick="return confirm('Are you sure you want to delete this video?')">Delete</a>
                </li>
            {{end}}
        </ul>
        <script>
        var tag = document.createElement('script');

        tag.src = "https://www.youtube.com/iframe_api";
        var firstScriptTag = document.getElementsByTagName('script')[0];
        firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);
        var player;
        function onPlayerStateChange(event) {
            if (event.data == YT.PlayerState.PLAYING ) {
                event.target.setPlaybackRate(2);
                event.target.setPlaybackQuality("hd1080");
                }
            }
            window.onYouTubeIframeAPIReady=function() {
                    player = new YT.Player('playerdiv', {
                        height: '80%',
                        width: '90%',
                        videoId: '{{index $.Links 0}}',
                        events: {
                            'onStateChange': onPlayerStateChange
                        }
                        });
                    }
        </script>

    </body>
</html>