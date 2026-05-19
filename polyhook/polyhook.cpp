#include "polyhook.h"

#include "polyhook2/Detour/x64Detour.hpp"
#include "polyhook2/ErrorLog.hpp"

#include <memory>
#include <mutex>
#include <string>

namespace {

class CaptureLogger : public PLH::Logger {
public:
    void log(const std::string& msg, PLH::ErrorLevel level) override {
        if (level >= PLH::ErrorLevel::WARN) {
            std::lock_guard<std::mutex> lk(m_mtx);
            m_last = msg;
        }
    }

    std::string pop() {
        std::lock_guard<std::mutex> lk(m_mtx);
        std::string msg;
        msg.swap(m_last);
        return msg;
    }

private:
    std::mutex  m_mtx;
    std::string m_last;
};

static std::shared_ptr<CaptureLogger> g_logger;

static void ensure_logger() {
    static std::once_flag once;
    std::call_once(once, []() {
        g_logger = std::make_shared<CaptureLogger>();
        PLH::Log::registerLogger(g_logger);
    });
}

} // namespace

struct PLH_Detour {
    PLH::x64Detour* detour     = nullptr;
    uint64_t        trampoline = 0;
    std::string     last_error;
};


extern "C" {

PLH_Detour* PLH_x64Detour_new(uint64_t fnAddress, uint64_t fnCallback) {
    ensure_logger();
    auto* d = new (std::nothrow) PLH_Detour();
    if (!d) return nullptr;
    d->detour = new (std::nothrow) PLH::x64Detour(fnAddress, fnCallback, &d->trampoline);
    if (!d->detour) {
        delete d;
        return nullptr;
    }
    return d;
}

int PLH_x64Detour_hook(PLH_Detour* d) {
    if (!d || !d->detour) return 0;
    g_logger->pop();
    const bool ok = d->detour->hook();
    if (!ok) d->last_error = g_logger->pop();
    return ok ? 1 : 0;
}

int PLH_x64Detour_unhook(PLH_Detour* d) {
    if (!d || !d->detour) return 0;
    g_logger->pop();
    const bool ok = d->detour->unHook();
    if (!ok) d->last_error = g_logger->pop();
    return ok ? 1 : 0;
}

uint64_t PLH_x64Detour_trampoline(PLH_Detour* d) {
    if (!d) return 0;
    return d->trampoline;
}

const char* PLH_x64Detour_last_error(PLH_Detour* d) {
    if (!d) return "";
    return d->last_error.c_str();
}

void PLH_x64Detour_free(PLH_Detour* d) {
    if (!d) return;
    delete d->detour;
    delete d;
}

} // extern "C"
