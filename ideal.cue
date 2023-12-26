package main
import (
    "github.com/DolusMockServer/dolus-expectations/dolus"
)

dolus.#Expectations & {
    expectations: [
        dolus.#Expectation & {
            priority: 1
            // add an expectation ID that will be generated
            request: dolus.#Request & {
                path:  "/store/order/{orderId}/p"
              method: "GET" 
            }
            response:  dolus.#Response & {
                body: {
                    petId: {
                        id: dolus.#GenInt32 &{min: 80, max:100}
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
