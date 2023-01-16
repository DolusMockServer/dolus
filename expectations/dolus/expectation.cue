package dolus


#Response: {
    body: _
    status: int
}

#Expectation:  {
    path: string
    method: "GET" | "POST"
    priority: int
    response:  #Response

}


#Expectations: {
  expectations:  [...#Expectation]
}

