package main
import (
	"github.com/MartinSimango/dolus/expectations/dolus"
)

dolus.#Expectations & {
    expectations: [
        dolus.#Expectation & {
        path:  "/store/order/2"
        method: "GET" 
        priority: -1
        response:  dolus.#Response & {
            body: {
				test: 5.2
				age: dolus.#GenInt32 
                shipDate: "a" 
                pair: 
                    g: {
                        hello: 20
                    }
                
			}
            status: 200
            }
        }
    ]
}
