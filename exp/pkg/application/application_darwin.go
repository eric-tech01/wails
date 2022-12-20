//go:build darwin

package application

/*

#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include "application.h"
#include "app_delegate.h"
#include "window_delegate.h"
#include <stdlib.h>

#import <Cocoa/Cocoa.h>

static AppDelegate *appDelegate = nil;

static void init(void) {
    [NSApplication sharedApplication];
    appDelegate = [[AppDelegate alloc] init];
    [NSApp setDelegate:appDelegate];

	[NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskLeftMouseDown handler:^NSEvent * _Nullable(NSEvent * _Nonnull event) {
		NSWindow* eventWindow = [event window];
		if (![eventWindow respondsToSelector:@selector(handleLeftMouseDown)]) {
			return event;
		}

		WindowDelegate* windowDelegate = (WindowDelegate*)[eventWindow delegate];
		[windowDelegate handleLeftMouseDown:event];
		return event;
	}];

	[NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskLeftMouseUp handler:^NSEvent * _Nullable(NSEvent * _Nonnull event) {
		NSWindow* eventWindow = [event window];
		if (![eventWindow respondsToSelector:@selector(handleLeftMouseUp)]) {
			return event;
		}

		WindowDelegate* windowDelegate = (WindowDelegate*)[eventWindow delegate];
		[windowDelegate handleLeftMouseUp:eventWindow];
		return event;
	}];
}

static void setActivationPolicy(int policy) {
    [NSApp setActivationPolicy:policy];
}

static void activateIgnoringOtherApps() {
	[NSApp activateIgnoringOtherApps:YES];
}

static void run(void) {
    @autoreleasepool {
        [NSApp run];
        [appDelegate release];
    }
}

// Destroy application
static void destroyApp(void) {
	[NSApp terminate:nil];
}

// Set the application menu
static void setApplicationMenu(void *menu) {
	NSMenu *nsMenu = (__bridge NSMenu *)menu;
	[NSApp setMainMenu:menu];
}

// Get the application name
static char* getAppName(void) {
	NSString *appName = [NSRunningApplication currentApplication].localizedName;
	if( appName == nil ) {
		appName = [[NSProcessInfo processInfo] processName];
	}
	return strdup([appName UTF8String]);
}

// get the current window ID
static unsigned int getCurrentWindowID(void) {
	NSWindow *window = [NSApp keyWindow];
	// Get the window delegate
	WindowDelegate *delegate = (WindowDelegate*)[window delegate];
	return delegate.windowId;
}

// Set the application icon
static void setApplicationIcon(void *icon, int length) {
    // On main thread
	dispatch_async(dispatch_get_main_queue(), ^{
		NSImage *image = [[NSImage alloc] initWithData:[NSData dataWithBytes:icon length:length]];
		[NSApp setApplicationIconImage:image];
	});
}

*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/exp/pkg/events"

	"github.com/wailsapp/wails/exp/pkg/options"
)

type macosApp struct {
	applicationMenu unsafe.Pointer
	parent          *App
}

func (m *macosApp) setIcon(icon []byte) {
	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *macosApp) name() string {
	appName := C.getAppName()
	defer C.free(unsafe.Pointer(appName))
	return C.GoString(appName)
}

func (m *macosApp) getCurrentWindowID() uint {
	return uint(C.getCurrentWindowID())
}

func (m *macosApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu for mac
		menu = defaultApplicationMenu()
	}
	menu.Update()

	// Convert impl to macosMenu object
	m.applicationMenu = (menu.impl).(*macosMenu).nsMenu
	C.setApplicationMenu(m.applicationMenu)
}

func (m *macosApp) run() error {
	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
		if m.parent.options != nil && m.parent.options.Mac != nil {
			C.setActivationPolicy(C.int(m.parent.options.Mac.ActivationPolicy))
		}
		C.activateIgnoringOtherApps()
	})
	C.run()
	return nil
}

func (m *macosApp) destroy() {
	C.destroyApp()
}

func newPlatformApp(app *App) *macosApp {
	appOptions := app.options
	if appOptions == nil {
		appOptions = options.ApplicationDefaults
	}
	C.init()
	return &macosApp{
		parent: app,
	}
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
	applicationEvents <- uint(eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &WindowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export processMessage
func processMessage(windowID C.uint, message *C.char) {
	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  C.GoString(message),
	}
}

//export processMenuItemClick
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}
