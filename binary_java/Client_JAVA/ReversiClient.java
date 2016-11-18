/*
 * ReversiClient v1.0
 * Yizhou Lin
 * ylin8 at Stevens.edu
 * 
 * Last Update: Sep 23 2012
 * 
 * This is a very simple program that shows how the client connects to the server.
 * 
 * It mainly has two parts/threads.
 * SendThread is a thread that read everything from standard input and then write send them to the socket; and
 * ReceiveThread is a thread that read everything from the socket and then write them to the standard output.
 */

import java.io.*;
import java.net.Socket;
public class ReversiClient
{
	public static String IP;
	public static int port;
	public static void main(String[] args)
	{
		if(args.length < 2)
		{
			System.out.println("usage: ReversiClient address port");
			//example:
			//java ReversiClient 127.0.0.1 4000
			return;
		}
		
		//get IP and port.
		IP = new String(args[0]);
		port = Integer.parseInt(args[1]);
		
		
		Socket s;
		try
		{
			s = new Socket(IP,port);
		
			//socket input
			BufferedReader in = new BufferedReader(new InputStreamReader(s.getInputStream()));
			
			//socket output
		   	PrintWriter out = new PrintWriter(s.getOutputStream());
		   	
		   	//standard input
			BufferedReader sysin = new BufferedReader(new InputStreamReader(System.in));   
			
			//establish two threads
			Thread send = new Thread(new SendThread(sysin,out));
			Thread receive = new Thread(new ReceiveThread(in));
			
			send.start();
			receive.start();
		}
		catch(Exception e)
		{
			e.printStackTrace();
		}
	}
}

class SendThread implements Runnable
{
	public BufferedReader sysin;
	PrintWriter out;
	public SendThread(BufferedReader si,PrintWriter o)
	{
		sysin = si;
		out = o;
	}
	
	public void run()
	{
		try
		{
			while(true)
			{
				String str=sysin.readLine();   
				out.println(str);   
				out.flush();   
			}
		}
		catch(Exception e)
		{
			e.printStackTrace();
		}
	}
}

class ReceiveThread implements Runnable
{
	public BufferedReader in;
	public ReceiveThread(BufferedReader br)
	{
		in=br;
	}
	public void run()
	{
		try
		{
			while(true)
			{
				System.out.println(in.readLine()); 
			}
		}
		catch(Exception e)
		{
			e.printStackTrace();
		}
	}
}