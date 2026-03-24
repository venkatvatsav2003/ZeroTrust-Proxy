import jwt
import time

SECRET = "super-secret-key-change-in-production"

def generate_token(sub, role):
    payload = {
        "sub": sub,
        "role": role,
        "iat": int(time.time()),
        "exp": int(time.time()) + 3600
    }
    return jwt.encode(payload, SECRET, algorithm="HS256")

if __name__ == "__main__":
    admin_token = generate_token("venkat", "admin")
    guest_token = generate_token("guest_user", "viewer")
    
    print("--- Zero-Trust Token Generator ---")
    print(f"\n[ADMIN TOKEN]:\n{admin_token}")
    print(f"\n[GUEST TOKEN (Blocked)]: \n{guest_token}")
    print("\n[USAGE EXAMPLE]:")
    print(f"curl -H 'Authorization: Bearer <TOKEN>' http://localhost:8443/api/data")
