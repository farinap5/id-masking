# IDOR mitigation example

This project contains a simple IDOR mitigation approach based on masks and the identity attribution directly to the ID passed as parameter. The idea is to hide the object ID veiling its incremental property.

This project tries to address some simple problems regarding the IDOR mitigation.
- Make discovery unfeasible through ID increment in brute force attacks;
- Decrease the difficulty of implementing mitigation at the data level in already consolidated projects.

For this POC the nounce could be anything, like the username, any number user specific or just the group the user is member of.

Before being sent to the frontend the ID is processed into a simple struct containing the users nounce and the object ID. The struct has both data separated by two dots `<nounce>:<id>`. This sequence is encrypted using a block cipher. Since the data could get too short asymmetric might be used as well.

The frontend must return this same encrypted message that will be reverted and verified before the ID being processed. The handler must verify the nounce before any action and forward the ID of successfully verified.

Id encoding implementation.
```lua
-- /person/list/
function handlerListPerson(*request, *response) {
	personList = listPerson(header['Token'].user)
	newList = []
	for i in personList {
    	    newList[i].id = base64(
        	    encrypt(personList[i].id + ":" + header['Token'].nounce)
    	    )
    	    newList[i].Name = personList[i].name
	}
	response.write(json(personList))
}
```

Id decoding implementation.
```lua
-- /person/get/
function handlerGetPerson(*request, *response) {
    nounce, personId = decode(query['id'])
    if header['Token'].nounce != nounce { response.write(error) return }
    
    person = getPersonDetails(personId)
    response.write(json(person))
}
```

The following code extracted from [pkg/server.go](pkg/server.go) contains the routes definition where you can see the implementation for all mechanism.
```go
	// init encryption function passing the key
	s := Sec{enc: *internal.Init("0123456789abcdef")}

	// IDOR vulnerable routes
	serverMux.HandleFunc("/person/list/", ListP)
	serverMux.HandleFunc("/person/get/", GetP)
	// Secured routes
	serverMux.HandleFunc("/person/list/secure/", s.ListPSecure)
	serverMux.HandleFunc("/person/get/secure/", s.GetPSecure)
```