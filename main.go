package main

import "fmt"

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	cl := NewPayNowApiClient(
		"b40245d2-ee08-4cf1-8396-7568b4e67306",
		"aace158f-5f8b-4250-a9b2-7844015a2a2f",
		"https://api.sandbox.paynow.pl/v3/",
	)
	resp, err := cl.CreatePayment("order123", 10000, "Test payment", "hkozacz@gmail.com")
	if err != nil {
		fmt.Println("Error creating payment:", err)
		return
	}
	fmt.Println(resp)

}
