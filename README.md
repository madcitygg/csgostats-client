csgostats-client
================
The csgostats-client is a simple log forwarder which will receive logs from a dedicated CS:GO server, and forward it to a webservice over http. It is small, efficient, easy to install and uses almost no resources.

Usage
-----
You'll want to install the csgostats-client on every server that you're running dedicated CS:GO servers on. This will improve reliability and make it easier to deal with network-wise.

**1. Download and unzip latest release**

Use the releases tab above to get the right release for your platform and architecture. Unzip it somewhere that's easy to find.

**2. Find your IP**

Use either `ifconfig` (if on Linux or OS X) or `ipconfig` if you're on Windows. The IP you want is the one of the main network card (usually called something like `en0`). This is necessary because SRCDS unfortunately doesn't seem to be able to send log messages over localhost or 127.0.0.1.

**3. Start csgostats-client**

You'll want to open a terminal (or cmd.exe on Windows), navigate to the folder where you saved `csgostats-client`. Then, start it with the port you want it to listen on:
```
./csgostats-client --port=26000
```
(If you want more output while testing the setup, add the -v flag, for verbose output)

**4. Configure SRCDS to send logs**

Add the following statements to the config file, or execute them at the SRCDS console. Be sure they are executed every time the server is restarted, or it will stop sending logs:
```
log on  // turns on logging
mp_logdetail 3  // log all attacks (so you can get ADR stats)
logaddress_add [IP]:[PORT]  // IP and port from step 2 and 3
```
Replace `[IP]` and `[PORT]` with the IP and port from step 2 and 3.

**5. Test**

Jump on the server and shoot some bots. If it's working, and you have verbose output on, you should start seeing a lot of log lines scroll by.


Architecture
------------
```
+------------------------------------+
|  Server 1                          |
|  +-----------+  UDP  +----------+  |
|  |  SRCDS 1  | +---> |          |  |          +----------------------+
|  +-----------+       |          |  |          |                      |
|  +-----------+  UDP  |          |  |   HTTP   |                      |
|  |  SRCDS 2  | +---> |  Client  +-----------> |     Stats Server     |
|  +-----------+       |          |  |          |                      |
|  +-----------+  UDP  |          |  |          |                      |
|  |  SRCDS 3  | +---> |          |  |          +----------------------+
|  +-----------+       +----------+  |
+------------------------------------+
```
