package main
import (
    "github.com/MartinSimango/dolus-expectations/dolus"
)

dolus.#Expectations & {
    expectations: [
        dolus.#Expectation & {
            request: dolus.#Request & {
                path:  "/store/order/2"
                method: "GET" 
            }
            response:  dolus.#Response & {
                body: {
                    age: dolus.#GenInt32 
                    name: "John Doe"        
                }
                status: 200
            }
      
        }
        
    ]
}
