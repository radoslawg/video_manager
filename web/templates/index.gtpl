<html>
    <head>
        {{template "style.gtpl"}}
    </head>
    <body>
        <h1>Files</h1>
        <ul>
        {{range .}}
            <li><a href="/view/{{.}}">{{.}}</a></li>
        {{end}}
        </ul>
    </body>
</html>