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
                    petId: {
                        id: "pete"
                    }
                    age: dolus.#GenInt32 
                    name: "John Doe"    
                    complete: "apple"    
                }
                status: 200
            }
      
        }
        
    ]
}
