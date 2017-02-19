<html>
    <head>
        <title>beego index template</title>
    </head>
    <body>
        <h1>Hello, world!{{.Username}}, {{.Email}}</h1>
        <div>
            {{if .TrueCond}}
            true condition.
            {{end}}
        </div>
        <div>
            {{if .FalseCond}}
            {{else}}
            false condition.
            {{end}}
        </div>
        <div>
            with Output:
            {{range .Nums}}
            {{.}}
            {{end}};
        </div>

        <div>
            range output:
            {{with .User}}
            Name:{{.Name}}; Age:{{.Age}}; Sex:{{.Sex}}
            {{end}}
        </div>

        <div>
            {{$tplVar := .TplVar}}
            The template variable is : {{$tplVar}}
        </div>

        <div>
            the result of template function is : {{htmlquote .Html}}
        </div>

        <div>
            pipeline : {{.Pipe | htmlquote}}
        </div>

        <div>
            {{template "netbsd"}}
        </div>
    </body>
</html>

{{define "netbsd"}}
Netbsd template test
{{end}}
