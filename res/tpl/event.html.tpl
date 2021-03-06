<!DOCTYPE html>
<html>

<head>
    <title>Event Registry - Event {{.Name}}</title>
    <link rel="icon" type="image/png" href="../assets/icon.png" />
    <link rel="stylesheet" href="../assets/bootstrap.min.css" />
</head>

<body>
    <div class="container">
        <h1>Event {{.Name}}</h1>

        <table class="table table-striped table-hover">
            <tr>
                <th>Slot</th>
                <th>Registrations</th>
                <th>Actions</th>
            </tr>
            {{- range $slotIndex, $slot := .Slots}}
            <tr>
                <td>{{.Name}}</td>
                <td class="text-{{if .IsFull}}danger{{else}}success{{end}}">{{len .Attendees}} / {{.Limit}}</td>
                <td>
                    {{- if not .IsFull}}
                    <a href="{{$.RegisterToken}}/register/{{$slotIndex}}">register</a>
                    {{- end}}
                </td>
            </tr>
            {{- end}}
        </table>

        <hr />
        <div class="text-center">
            <img src="../assets/logo.png" title="Logo" alt="Logo" />
            <br />
            <a href="../assets/license.txt">License</a>
        </div>
    </div>
</body>

</html>
