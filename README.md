# gomemorycache

# Install
## go get -u github.com/v-kolodii/gomemorycache

## Usage:
    func main() {
        cache := gomemorycache.New()

        cache.Set("userId", 42)
        userId := cache.Get("userId")

        fmt.Println(userId)

        cache.Delete("userId")
        userId := cache.Get("userId")

        fmt.Println(userId)
    }