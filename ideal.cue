package main
import (
    "github.com/DolusMockServer/dolus/cue-expectations/core"
    "github.com/DolusMockServer/dolus/cue-expectations/core/task"
)

core.#Expectations & {
    expectations: [
        core.#Expectation & {
            priority: 1
            // add an expectation ID that will be generated
            request: core.#Request & {
                path:  "/store/order/2/p"
              method: "GET" 
            }
            response:  core.#Response & {
                body: {
                    petId: {
                        id: task.#GenInt32 &{min: 80, max:100}
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
