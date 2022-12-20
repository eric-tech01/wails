//go:build darwin

#ifndef WindowDelegate_h
#define WindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler, WKNavigationDelegate>

@property bool hideOnClose;
@property (retain) WKWebView* webView;
@property unsigned int windowId;
@property (retain) NSEvent* leftMouseEvent;
@property unsigned int invisibleTitleBarHeight;

- (void)handleLeftMouseUp:(NSWindow *)window;
- (void)handleLeftMouseDown:(NSEvent*)event;

@end


#endif /* WindowDelegate_h */
