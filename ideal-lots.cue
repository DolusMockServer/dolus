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
				path:   "/store/order/3/p?value=3&age=5" // TODO: if path param is left out it still works i.e /store/order//p still works
				method: "GET"
		
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
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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
		expectation.#Expectation & {
			priority: 1
			// add an expectation ID that will be generated
			request: expectation.#Request & {
				path:   "/store/order/3/p"
				method: "GET"
				headers: {
					"Content-Type": ["application/json"]
					"T" : ["3","4"]
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

