<!DOCTYPE html>
<html>

<head>
    <title>Event Registry - Event {{.Name}}</title>
    <link rel="icon" type="image/png" href="../../assets/icon.png" />
    <link rel="stylesheet" href="../../assets/bootstrap.min.css" />
</head>

<body>
    <div class="container">
        <h1>Event {{.Name}}</h1>

        <table class="table table-striped table-hover">
            <tr>
                <th>Slot</th>
                <th>Registrations</th>
                <th>Attendees</th>
            </tr>
            {{- range $slotIndex, $slot := .Slots}}
            <tr>
                <td>{{$slot.Name}}</td>
                <td class="text-{{if .IsFull}}danger{{else}}success{{end}}">{{len .Attendees}} / {{.Limit}}</td>
                <td>
                    <ol>
                        {{- range $attendee := $slot.Attendees}}
                        <li>{{$attendee}} <a href="admin/{{$slotIndex}}/delete/{{$attendee}}">&#10060;</a></li>
                        {{- end}}
                    </ol>
                </td>
            </tr>
            {{- end}}
        </table>

        <hr />
        <div class="text-center">
            <img src="../../assets/logo.png" title="Logo" alt="Logo" />
            <br />
            <a href="../../assets/license.txt">License</a>
        </div>
    </div>
</body>

</html>
