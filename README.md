# UE4-Nakama-matchmaking-backend-interface
This minimalistic go plugin example for nakama shows how you can handle matchmaking with nakama and unreal engine default net code .
#How to use 

-install this plugin for sending jsons https://github.com/Stefander/JSONQuery.

-install nakama  https://github.com/heroiclabs/nakama-unreal authenticate and make a rtclient .

-Create a minimal go plugin https://github.com/heroiclabs/nakama/tree/master/sample_go_module.

-paste  this go code into your plugin code.

-implement this blueprints inside your game (this should be called only by dedicated servers . better be placed in GameMode)
https://blueprintue.com/blueprint/epzpnat1/ .




#How it works?
1. blueprint calls a rpc on server and passes server port to nakama .
2.server ip is submitted in a array .
3. when match is made nakama returns port and ip as a string instead of match id (plugin only passes servers which were pinged less than a  constant time .


Getportfunction :
```sh
int ATankmode::GetPort()
{
	int output;
	if (GetWorld())
		output = GetWorld()->URL.Port;
	else output = -1;
	return output;
}

```
