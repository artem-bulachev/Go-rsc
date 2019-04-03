// forms.go
package main

import (
    "fmt"
	"net/http"
    "strconv"
)

type eData struct {
    owner string
    amount int
}

type Event struct{
    eType string
    data eData
}

type Account struct {
    id string
    eCnt int
    events [100]Event
}

var bank [1000]Account
var idCnt int

func main() {

    idCnt = 1

	http.HandleFunc("/CreateAcc", CreateAccount)
	http.HandleFunc("/UpdateOwner", UpdateOwner)
	http.HandleFunc("/Deposit", Deposit)
	http.HandleFunc("/Withdrawal", Withdrawal)

	http.ListenAndServe(":8080", nil)
}

// http://localhost:8080/CreateAcc?owner="sunny"
func CreateAccount(w http.ResponseWriter, r *http.Request) {
    var newAcc Account
    newAcc.id = strconv.Itoa(idCnt)
    owner := r.URL.Query().Get("owner")
    if owner == "" {
        w.Write([]byte("Error!"))
        return
	}
    var detail eData
    detail.owner = owner
    var newEvent Event
    newEvent.eType = "CreateAccount"
    newEvent.data = detail
    newAcc.events[0] = newEvent
    newAcc.eCnt = 1;

    bank[idCnt-1] = newAcc
    idCnt = idCnt + 1
    w.Write([]byte("Created at ID: " + newAcc.id))
    fmt.Println(newAcc)

    var status eData
    status = getStatus(bank[idCnt-2])
    w.Write([]byte(" Status:(owner:" + status.owner + " amount:" + strconv.Itoa(status.amount) + ")"))
}

// http://localhost:8080/UpdateOwner?id=1&owner='sunny1'
func UpdateOwner(w http.ResponseWriter, r *http.Request) {
    cIDstr := r.URL.Query().Get("id")
    cID, err := strconv.Atoi(cIDstr)
    if err != nil{
        w.Write([]byte("Error! Can't get ID with " + cIDstr + " " + strconv.Itoa(cID)))
        return
    }
    cOwner := r.URL.Query().Get("owner")
    if cOwner == "" {
        w.Write([]byte("Error! Can't get owner!"))
        return
    }

    if idCnt <= cID{
        w.Write([]byte("non Existing ID: " + cIDstr))
        return
    }

    var detail eData
    detail.owner = cOwner
    var newEvent Event
    newEvent.eType = "UpdateOwner"
    newEvent.data = detail
    bank[cID-1].events[bank[cID-1].eCnt] = newEvent
    bank[cID-1].eCnt = bank[cID-1].eCnt + 1;

    w.Write([]byte("Updated owner: " +cOwner + " with ID " + cIDstr))
    fmt.Println(bank[cID-1])

    var status eData
    status = getStatus(bank[cID-1])
    w.Write([]byte(" Status:(owner:" + status.owner + " amount:" + strconv.Itoa(status.amount) + ")"))
}

// http://localhost:8080/Deposit?id=1&amount=12
func Deposit(w http.ResponseWriter, r *http.Request) {
    cIDstr := r.URL.Query().Get("id")
    cID, err := strconv.Atoi(cIDstr)
    if err != nil {
        w.Write([]byte("Error! Can't get ID with" + cIDstr + strconv.Itoa(cID)))
        return
    }
    cAmountstr := r.URL.Query().Get("amount")
    cAmount, err := strconv.Atoi(cAmountstr)
    if err != nil {
        w.Write([]byte("Error! Can't get amount!"))
        return
    }

    if idCnt <= cID{
        w.Write([]byte("non Existing ID: " + cIDstr))
        return
    }

    var detail eData
    detail.amount = cAmount
    var newEvent Event
    newEvent.eType = "Deposit"
    newEvent.data = detail
    bank[cID-1].events[bank[cID-1].eCnt] = newEvent
    bank[cID-1].eCnt = bank[cID-1].eCnt + 1;

    w.Write([]byte("Deposit: " +cAmountstr + " with ID " + cIDstr))
    fmt.Println(bank[cID-1])

    var status eData
    status = getStatus(bank[cID-1])
    w.Write([]byte(" Status:(owner:" + status.owner + " amount:" + strconv.Itoa(status.amount) + ")"))
}

// http://localhost:8080/Withdrawal?id=1&amount=12
func Withdrawal(w http.ResponseWriter, r *http.Request) {
    cIDstr := r.URL.Query().Get("id")
    cID, err := strconv.Atoi(cIDstr)
    if err != nil {
        w.Write([]byte("Error! Can't get ID with" + cIDstr + strconv.Itoa(cID)))
        return
    }
    cAmountstr := r.URL.Query().Get("amount")
    cAmount, err := strconv.Atoi(cAmountstr)
    if err != nil {
        w.Write([]byte("Error! Can't get amount!"))
        return
    }

    if idCnt <= cID{
        w.Write([]byte("non Existing ID: " + cIDstr))
        return
    }

    var detail eData
    detail.amount = cAmount
    var newEvent Event
    newEvent.eType = "Withdrawal"
    newEvent.data = detail
    bank[cID-1].events[bank[cID-1].eCnt] = newEvent
    bank[cID-1].eCnt = bank[cID-1].eCnt + 1;

    w.Write([]byte("Withdrawal: " +cAmountstr + " with ID " + cIDstr))
    fmt.Println(bank[cID-1])

    var status eData
    status = getStatus(bank[cID-1])
    w.Write([]byte(" Status:(owner:" + status.owner + " amount:" + strconv.Itoa(status.amount) + ")"))

}

func getStatus(acc Account) eData{
    var status eData
    for i := 0; i < acc.eCnt ; i++ {
        e := acc.events[i]
        if e.eType == "CreateAccount" {
            status.owner = e.data.owner
            status.amount = 0
        }
        if e.eType == "UpdateOwner" {
            status.owner = e.data.owner
        }
        if e.eType == "Deposit" {
            status.amount = status.amount + e.data.amount
        }
        if e.eType == "Withdrawal" {
            status.amount = status.amount - e.data.amount
        }
	}
    return status
}