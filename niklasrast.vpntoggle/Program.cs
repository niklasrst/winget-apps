using System;
using System.Diagnostics;
using System.Linq;
using System.ServiceProcess;
using System.ServiceProcess; // Ensure this using directive is present for ServiceControllerStatus

// Make sure to add a reference to System.ServiceProcess.ServiceController in your project file if targeting .NET Core/5+/6+
// For example, run: dotnet add package System.ServiceProcess.ServiceController

class Program
{
    static void Main(string[] args)
    {
        // /help parameter
        if (args.Length > 0 && args[0] == "/help")
        {
            Console.WriteLine("This tool can toggle VPN on or off");
            Console.WriteLine("Usage: vpntoggle");
            return;
        }

        // Find VPN services
        var vpnServices = ServiceController.GetServices()
            .Where(s => s.ServiceName.IndexOf("vpn", StringComparison.OrdinalIgnoreCase) >= 0)
            .ToList();

        if (!vpnServices.Any())
        {
            Console.WriteLine("No VPN service to toggle.");
            Environment.Exit(0);
        }

        // Check if any VPN service is running
        bool anyRunning = vpnServices.Any(s => s.Status == ServiceControllerStatus.Running);

        if (!anyRunning)
        {
            // Start VPN services
            bool allStarted = true;
            foreach (var svc in vpnServices)
            {
                try
                {
                    if (svc.Status != ServiceControllerStatus.Running)
                    {
                        svc.Start();
                        svc.WaitForStatus(ServiceControllerStatus.Running, TimeSpan.FromSeconds(10));
                    }
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error starting service {svc.ServiceName}: {ex.Message}");
                    allStarted = false;
                }
            }

            if (allStarted)
            {
                Console.WriteLine("VPN services started successfully.");

                // Launch VPN client
                string vpnClientPath = @"C:\Program Files (x86)\Cisco\Cisco AnyConnect Secure Mobility Client\vpnui.exe";
                try
                {
                    Process.Start(vpnClientPath);
                    Console.WriteLine("VPN Client launched successfully.");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error launching VPN client: {ex.Message}");
                }
            }
            else
            {
                Console.WriteLine("Error starting one or more services.");
            }
        }
        else
        {
            // Stop VPN services
            bool allStopped = true;
            foreach (var svc in vpnServices)
            {
                try
                {
                    if (svc.Status == ServiceControllerStatus.Running)
                    {
                        svc.Stop();
                        svc.WaitForStatus(ServiceControllerStatus.Stopped, TimeSpan.FromSeconds(10));
                    }
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error stopping service {svc.ServiceName}: {ex.Message}");
                    allStopped = false;
                }
            }

            if (allStopped)
            {
                Console.WriteLine("VPN services stopped successfully.");
            }
            else
            {
                Console.WriteLine("Error stopping one or more services.");
            }

            // Stop VPN processes
            var vpnProcesses = Process.GetProcesses()
                .Where(p => p.ProcessName.IndexOf("vpn", StringComparison.OrdinalIgnoreCase) >= 0)
                .ToList();

            bool allProcStopped = true;
            foreach (var proc in vpnProcesses)
            {
                try
                {
                    proc.Kill();
                    proc.WaitForExit(5000);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error stopping process {proc.ProcessName}: {ex.Message}");
                    allProcStopped = false;
                }
            }

            if (allProcStopped)
            {
                Console.WriteLine("VPN processes stopped successfully.");
            }
            else
            {
                Console.WriteLine("Error stopping one or more processes.");
            }
        }
    }
}
