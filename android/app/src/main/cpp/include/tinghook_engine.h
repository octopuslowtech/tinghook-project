#ifndef TINGHOOK_ENGINE_H
#define TINGHOOK_ENGINE_H

namespace tinghook {

class Engine {
public:
    static void init();
    static void connect(const char* apiKey, const char* deviceUid);
    static void disconnect();
};

}

#endif // TINGHOOK_ENGINE_H
