Event Registry

1. run evreg.exe
2. go to http://localhost:8080/event/R3G1S7ER
3. click any register link
4. enter name and hit submit
5. go back to http://localhost:8080/event/R3G1S7ER
6. observe one new registration
7. go to http://localhost:8080/event/L3TM3IN/admin
8. observe list of attendees

- register.json is read on application startup
- register.json is written during application runtime, do not touch!
- link with RegisterToken (here R3G1S7ER) is given to students
- link with AdminToken (here L3TM3IN) is given to staff
- attendees can be removed via admin view with one click
- interface / port can be specified as first commandline argument
  e.g. evreg.exe 192.168.1.42:80

Run:

    sudo apt install golang
    go run ./cmd/evreg
