"use client"
import {useEffect, useState} from "react";

const LogSteam = () =>{
    const [logs, setLogs] = useState<string[]>([]);
    const {socket, setSocket} = useState<WebSocket | null>(null);
    const [command, setCommand] = useState('');

    useEffect(() => {
        const ws = new WebSocket('ws://localhost:8080');
        ws.onopen = () => {
            console.log('WebSocket已经打开');
        };

        ws.onmessage = (event) =>{
            console.log('日志消息：', event.data);
            setLogs((prevLogs: any) => [...prevLogs, event.data]);
        }

        return() =>{
            ws.close();
        }
    }, []);

    return (
        <div className="container mx-auto p-4">
            <h1 className="text-2xl font-bold mb-4">服务器日志流</h1>

            <div className="bg-gray-100 p-4 rounded-md mb-4">
                <h2 className="text-xl font-semibold mb-2">Logs</h2>
                <div className="h-64 overflow-auto border border-gray-300 p-2 bg-white">
                    {logs.map((log, index) => (
                        <p key={index} className="text-sm text-gray-700">
                            {log}
                        </p>
                    ))}
                </div>
            </div>

            <div className="bg-gray-100 p-4 rounded-md">
                <h2 className="text-xl font-semibold mb-2">输入指令</h2>
                <input
                    type="text"
                    className="border p-2 w-full rounded-md mb-2"
                    value={command}
                    onChange={(e) => setCommand(e.target.value)}
                    placeholder="输入指令发送到后端...."
                />
            </div>
        </div>
    );
};

export default LogSteam;
