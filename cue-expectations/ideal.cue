package main

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/cue-expectations/core/task"
)

expectation.#Expectations & {
	expectations: [
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/{orderId}/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
				}
				params: {
					query: {
						value: ["3","7"]
						age: ["5"]	
					}
					path: {
						orderId: expectation.#Matcher & { match: "eq",	value: "123" }
					}
				}

			}

			// cookies: {
			//   "cookie1": core.#Matcher & { match: "eq",	value: "value1" }
			// }

			// 	}
			response: expectation.#Response & {
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

		},

	]
}