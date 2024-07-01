import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

public class HashUtils {
    public static String sha256(String val) {
        try {
            MessageDigest instance = MessageDigest.getInstance("SHA-256");
            byte[] encodedHash = instance.digest(val.getBytes());
            // 将字节数组转换为十六进制字符串
            StringBuilder hexString = new StringBuilder(2 * encodedHash.length);
            for (byte b : encodedHash) {
                String hex = Integer.toHexString(0xff & b);
                if (hex.length() == 1) {
                    hexString.append('0');
                }
                hexString.append(hex);
            }
            return hexString.toString();
        }
        catch(NoSuchAlgorithmException e) {
            System.out.println("impossible hhh");
        }
        return "";
    }
}   
