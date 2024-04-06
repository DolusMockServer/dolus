package main

import (
	"github.com/DolusMockServer/dolus/pkg/expectation"
	"github.com/DolusMockServer/dolus/cue-expectations/core/task"
)

expectation.#Expectations & {
	expectations: [
		expectation.#Expectation & {
			priority: 2
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/2/p?value=3&age=5"
				method: "GET"
				headers: {
					"Content-Type": "application/json"
					"T": ["3", "4"]
					"G": expectation.#Matcher & {match: "regex", value: "hello"}
				}
				// params: {
				// 	query: { 
				// 		value: [3]
				// 		age: expectation.#Matcher & {match: "eq", value: 2}
				// 	}
				// }
				cookies: [
					expectation.#Cookie & {name: "cookie1", value: "As", path: "/p"},
					// expectation.#CookieMatcher & { match: "eq", value: expectation.#Cookie &{name: "cookie2", value: "Bs", path: "/p" }}
				]

			}

			response: expectation.#Response & {
				body: {
					petId: {
						id: task.#GenInt32 & {min: 80, max: 100}
					}
					// age: dolus.#GenInt32 
					// name: "John Doe"    
					complete: true
					status:   "approved"
				}

				status: 200
			}

		},

		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T": ["3", "4"]
				}
				params: {
					query: {
						value: [3]
						age: ["5"]
					}
					path: {
						orderId: "3"
					}
				}

			}

			response: expectation.#Response & {
				body: {
					petId: {
						id: task.#GenInt32 & {min: 80, max: 100}
					}

					complete: false
					status:   "good day"
				}

				status: 200
			}

		},
	]
}
