package main

import (
	"github.com/DolusMockServer/dolus/pkg/expectation/cue"
	"github.com/DolusMockServer/dolus/cue-expectations/core/task"
)

cue.#Expectations & {
	expectations: [
		cue.#Expectation & {
		priority: 1
			// add an expectation ID that will be generated
		request: cue.#Request & {
			path:   "/store/order/2/p?param1&param2"
			method: "GET"	
			params: {
				query: {
					param1: cue.#Matcher & { match: "eq", value: "value1"} 
					param2: cue.#Matcher & { match: "has", value: "t"}
				}
				path: {
					param1: cue.#Matcher & { match: "eq" , value: "value1"}
					param2: cue.#Matcher & { match: "has" }
				}
			}
		}
		
        // headers: {
        //   "Content-Type": core.#Matcher & { match: "eq",	value: "application/json" }
        // }
   
        // cookies: {
        //   "cookie1": core.#Matcher & { match: "eq",	value: "value1" }
        // }

		// 	}
		response: cue.#Response & {
			body: {
				petId: {
					id: task.#GenInt32 & {min: 80, max: 100}
				}
				// age: dolus.#GenInt32 
				// name: "John Doe"    
				complete: true
				status:   "good day"
			}

			status: 200
		}

		}

	]
}
