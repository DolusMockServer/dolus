package main
import (
	"github.com/MartinSimango/dolus/expectations/dolus"
)

dolus.#Expectations & {
    expectations: [
        dolus.#Expectation & {
        path:  "/store/order/1"
        method: "GET" 
        priority: 1
        response:  dolus.#Response & {
            body: {
				test:  2
				age: dolus.#GenInt32 
			}
            status: 2
            }
        }
    ]
}
