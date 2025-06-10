using System;
using System.Diagnostics;
using Microsoft.Win32;

class Program
{
    static void Main()
    {
        try
        {
            // Open the registry key
            using (RegistryKey accountsKey = Registry.LocalMachine.OpenSubKey(@"SOFTWARE\Microsoft\Provisioning\OMADM\Accounts"))
            {
                if (accountsKey == null)
                {
                    Console.WriteLine("Registry path not found.");
                    return;
                }

                // Get subkey names (account identifiers)
                string[] accountNames = accountsKey.GetSubKeyNames();

                foreach (string account in accountNames)
                {
                    Console.WriteLine($"Enrolling account: {account}");
                    Process process = new Process();
                    process.StartInfo.FileName = "deviceenroller.exe";
                    process.StartInfo.Arguments = $"/o {account} /c /b";
                    process.StartInfo.UseShellExecute = false;
                    process.StartInfo.CreateNoWindow = true;

                    process.Start();
                    process.WaitForExit();
                }
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine("An error occurred:");
            Console.WriteLine(ex.Message);
        }
    }
}
