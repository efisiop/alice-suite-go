#!/usr/bin/env python3
import json
import requests
import time
import threading
import sys

def test_health_endpoint():
    time.sleep(1)
    try:
        response = requests.get('http://localhost:8080/api/health', timeout=2)
        print("=== Health API Test ===")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        return response.status_code == 200
    except Exception as e:
        print(f"Health API Error: {e}")
        return False

def test_user_endpoint():
    time.sleep(1)
    try:
        response = requests.get('http://localhost:8080/api/user?user_id=test123', timeout=2)
        print("=== User API Test ===")
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        return response.status_code in [200, 404]  # 404 is expected if user doesn't exist
    except Exception as e:
        print(f"User API Error: {e}")
        return False

if __name__ == '__main__':
    try:
        import subprocess
        server_proc = subprocess.Popen(['go', 'run', './cmd/reader'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        
        time.sleep(2)  # Wait for server to start

        # Test the endpoints
        health_ok = test_health_endpoint()
        user_ok = test_user_endpoint()

        if health_ok and user_ok:
            print("✅ All API tests passed!")
        else:
            print("❌ Some API tests failed")
            print("Server logs:")
            stdout, stderr = server_proc.communicate(timeout=1)
            if stdout:
                print("STDOUT:", stdout)
            if stderr:
                print("STDERR:", stderr)

    except KeyboardInterrupt:
        print("\nStopping server...")
    finally:
        server_proc.terminate()
        server_proc.wait()
