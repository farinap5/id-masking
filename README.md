This project contains a simple IDOR mitigation approach based on masks and the identity attribution directly to the ID passed as parameter. The idea is to hide the object ID veiling its incremental property.

This project tries to address some simple problems regarding the IDOR mitigation.
- Make discovery unfeasible through ID increment in brute force attacks;
- Decrease the difficulty of implementing mitigation at the data level in already consolidated projects.

For this POC the nounce could be anything, like the username, any number user specific or just the group the user is member of.

Before being sent to the frontend the ID is processed into a simple struct containing the users nounce and the object ID. The struct has both data separated by two dots `<nounce>:<id>`. This sequence is encrypted using a block cipher. Since the data could get too short asymmetric might be used as well.

The frontend must return this same encrypted message that will be reverted and verified before the ID being processed. The handler must verify the nounce before any action and forward the ID of successfully verified.


AES:ENC(`123:1`) = WXYZ=
AES:DEC(`WXYZ=`) = 123:1