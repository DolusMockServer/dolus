package main
import (
    "github.com/MartinSimango/dolus-expectations/dolus"
)

dolus.#Expectations & {
    expectations: [
        dolus.#Expectation & {
            // add an expectation ID that will be generated
            request: dolus.#Request & {
                path:  "/store/order/2"
                method: "GET" 
            }
            response:  dolus.#Response & {
                body: {
                    petId: {
                        id: dolus.#GenInt32 
                    }
                    // age: dolus.#GenInt32 
                    // name: "John Doe"    
                    complete: true   
                    status: "good day"
                }
                status: 200
            }
      
        }
        
    ]
}
