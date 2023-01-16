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
				test: 5.0 
				age: dolus.#GenInt32 & { min: 20 , max: 40}
			}
            status: 2
            }
        }
    ]
}
