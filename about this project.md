# Project Overview: Local Network Dashboard & Communication Tool

## What This Project Is (In Plain English)
You are building a local network dashboard + communication tool.  
**That’s it.** Not a research lab. Not a mega system.  

Just a clean, understandable LAN tool that teaches you networking.

---

## Core Purpose
See, understand, and communicate with devices on your local network.

---

## Final Project Functionalities (Clear + Minimal)

### 1️⃣ Show My Network Info
When I open the app, I should see:
* **My device name**
* **My local IP address**
* **My subnet mask**
* **My default gateway (router IP)**
> 👉 *This teaches how devices get network identity.*

### 2️⃣ Show All Devices on My LAN
The app should show:
* List of devices connected to my WiFi / LAN
* Their IP addresses
* Which ones are active
> 👉 *This teaches LAN scanning + network visibility.*

### 3️⃣ Discover Other App Users Automatically
If another device is running this same app on the same WiFi:
* They should automatically appear in my app
* I should see their **Device Name** and **IP**
> 👉 *This teaches peer discovery + broadcasting.*

### 4️⃣ Send Simple Messages Between Devices (LAN Only)
* **Example:** I type a message on my laptop; my phone (on the same WiFi) receives it instantly.
> 👉 *This teaches LAN communication + real-time networking.*

### 5️⃣ Share Files Over LAN
* **Example:** Send a file from laptop → phone with no internet involved.
> 👉 *This teaches local file transfer + HTTP streaming.*

---

## That Is The Entire Project
No more. No extra complexity. No buzzwords. Just:
* **See** your network
* **See** who is on it
* **Talk** to them
* **Share** files

---

## What You Are NOT Building (Right Now)
* ❌ Internet exposure
* ❌ NAT traversal
* ❌ VPNs
* ❌ VPS reverse proxies
* ❌ Advanced security models  
*These come later, optionally.*

---

## Simple Mental Model
**Your app = LAN WhatsApp + LAN Network Scanner**

### Simple User Flow
1.  Open app
2.  See my IP + network
3.  See all connected devices
4.  See devices running same app
5.  Click a device → send message or file

---

## Technical Breakdown
| Component | Technology | Responsibility |
| :--- | :--- | :--- |
| **Frontend** | React | UI dashboard, Device list, Chat + file UI |
| **Backend** | Go | Network info, LAN scanning, Peer discovery, Messaging + file APIs |

---

## Why This Project Is Perfect
It teaches the following without overwhelming you:
* IP addressing & Subnets
* NAT & LAN routing
* Broadcast / Multicast
* TCP / UDP
* HTTP streaming

---

## Final Clarification Check
You are building: **A LAN network dashboard + communication tool. Nothing more.**
