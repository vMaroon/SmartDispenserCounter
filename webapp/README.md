# Smart Dispenser Counter WebApp
Our frontend consists of a wix web-application that can be used on desktop/mobile comfortably.
The UI of wix is not provided in this repository, but the backend code is available.

The backend code is a collection of javascript modules that wrap the functionality required to communicate with and use
our RESTful API server.

---
## Mobile View
<img src="images/mobile.png" width="40%" height="40%">

---
## Desktop View
<img src="images/desktop.png" width="60%" height="60%">

---
Requests can be added through the page, new requests are automatically made "active" - meaning that 
the ESP32 controller would receive them for serving.

The history of the 50 latest requests can be seen in the page. A previous request cannot be made active again.