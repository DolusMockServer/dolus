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
				path:   "/store/order/2/p?param1&param2"
				method: "GET"
				queryParams: {
					param1: core.#Matcher & {match: "eqa", value: "value1"}
					param2: core.#Matcher & {match: "has"}
				}
				headers: {
					"Content-Type": core.#Matcher & {match: "eq", value: "application/json"}
				}
				cookies: {
					"cookie1": core.#Matcher & {match: "eq", value: "value1"}
				}

			}
			response: core.#Response & {
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
