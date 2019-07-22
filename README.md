# udp-chat

This is a simple example of a UDP chat client that utilizes multicast IP to deliver messages. 
Users can change their name with the `/name` command. 

Since this is meant to be a basic example, there is no persistent state of user settings, 
no authentication, no encryption, no timestamp displayed for messages, etc.


### Does this work across the internet?

No. UDP multicast is not routable across the public internet. Depending on local network configuration, 
it should work between computers on a home network or business/university.
