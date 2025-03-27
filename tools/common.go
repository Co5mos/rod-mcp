package tools

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod-mcp/types"
	"github.com/go-rod/rod-mcp/utils"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	defaultWaitStableDur = 1 * time.Second
	defaultDomDiff       = 0.2
)

var (
	Navigation = mcp.NewTool("rod_navigate",
		mcp.WithDescription("Navigate to a URL"),
		mcp.WithString("url", mcp.Description("URL to navigate to"), mcp.Required()),
	)
	GoBack = mcp.NewTool("rod_go_back",
		mcp.WithDescription("Go back in the browser history, go back to the previous page"),
	)
	GoForward = mcp.NewTool("rod_go_forward",
		mcp.WithDescription("Go forward in the browser history, go to the next page"),
	)
	ReLoad = mcp.NewTool("rod_reload",
		mcp.WithDescription("Reload the current page"),
	)
	PressKey = mcp.NewTool("rod_press_key",
		mcp.WithDescription("Press a key on the keyboard"),
		mcp.WithString("key", mcp.Description("Name of the key to press or a character to generate, such as `ArrowLeft` or `a`"), mcp.Required()),
	)
	Pdf = mcp.NewTool("rod_pdf",
		mcp.WithDescription("Generate a PDF from the current page"),
		mcp.WithString("file_path", mcp.Description("Path to save the PDF file"), mcp.Required()),
		mcp.WithString("file_name", mcp.Description("Name of the PDF file"), mcp.Required()),
	)
	CloseBrowser = mcp.NewTool("rod_close",
		mcp.WithDescription("Close the browser"),
	)
	Screenshot = mcp.NewTool("rod_screenshot",
		mcp.WithDescription("Take a screenshot of the current page or a specific element"),
		mcp.WithString("name", mcp.Description("Name of the screenshot"), mcp.Required()),
		mcp.WithString("selector", mcp.Description("CSS selector of the element to take a screenshot of")),
		mcp.WithNumber("width", mcp.Description("Width in pixels (default: 800)")),
		mcp.WithNumber("height", mcp.Description("Height in pixels (default: 600)")),
	)
	Click = mcp.NewTool("rod_click",
		mcp.WithDescription("Click an element on the page"),
		mcp.WithString("selector", mcp.Description("CSS selector of the element to click"), mcp.Required()),
	)
	Fill = mcp.NewTool("rod_fill",
		mcp.WithDescription("Fill out an input field"),
		mcp.WithString("selector", mcp.Description("CSS selector of the element to type into"), mcp.Required()),
		mcp.WithString("value", mcp.Description("Value to fill"), mcp.Required()),
	)
	Selector = mcp.NewTool("rod_selector",
		mcp.WithDescription("Select an element on the page with Select tag"),
		mcp.WithString("selector", mcp.Description("CSS selector for element to select"), mcp.Required()),
		mcp.WithString("value", mcp.Description("Value to select"), mcp.Required()),
	)
	Evaluate = mcp.NewTool("rod_evaluate",
		mcp.WithDescription("Execute JavaScript in the browser console"),
		mcp.WithString("script", mcp.Description("JavaScript code to execute"), mcp.Required()),
	)
	Drag = mcp.NewTool("rod_drag",
		mcp.WithDescription("Perform drag and drop between two elements"),
		mcp.WithString("source_selector", mcp.Description("CSS selector of the source element"), mcp.Required()),
		mcp.WithString("target_selector", mcp.Description("CSS selector of the target element"), mcp.Required()),
	)
	Wait = mcp.NewTool("rod_wait",
		mcp.WithDescription("Wait for a specified time in seconds"),
		mcp.WithNumber("time", mcp.Description("The time to wait in seconds (capped at 10 seconds)"), mcp.Required()),
	)
	Snapshot = mcp.NewTool("rod_snapshot",
		mcp.WithDescription("Capture accessibility snapshot of the current page"),
	)
	SelectElement = mcp.NewTool("rod_select_element",
		mcp.WithDescription("Select and retrieve information about elements on the page"),
		mcp.WithString("selector", mcp.Description("CSS selector of the elements to retrieve"), mcp.Required()),
		mcp.WithBoolean("get_attributes", mcp.Description("Whether to include element attributes in the result (default: false)")),
		mcp.WithBoolean("get_text", mcp.Description("Whether to include element text content in the result (default: true)")),
	)
	GetText = mcp.NewTool("rod_get_text",
		mcp.WithDescription("Get text content from the page or a specific element"),
		mcp.WithString("selector", mcp.Description("CSS selector of the element to get text from (if empty, gets text from entire page)")),
		mcp.WithBoolean("trim", mcp.Description("Whether to trim whitespace from the text (default: true)")),
		mcp.WithBoolean("include_hidden", mcp.Description("Whether to include hidden text (default: false)")),
	)
)

type ToolHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)

var (
	NavigationHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			url := request.Params.Arguments["url"].(string)
			if !utils.IsHttp(url) {
				return nil, errors.New("invalid URL")
			}

			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to navigate to %s: %s", url, err.Error()))
			}
			err = page.Navigate(url)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to navigate to %s: %s", url, err.Error()))
			}
			page.WaitDOMStable(defaultWaitStableDur, defaultDomDiff)
			return mcp.NewToolResultText(fmt.Sprintf("Navigated to %s", url)), nil
		}
	}

	GoBackHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to go back: %s", err.Error()))
			}
			err = page.NavigateBack()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to go back: %s", err.Error()))
			}
			page.WaitDOMStable(defaultWaitStableDur, defaultDomDiff)
			return mcp.NewToolResultText("Go back successfully"), nil
		}
	}

	GoForwardHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to go forward: %s", err.Error()))
			}
			err = page.NavigateForward()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to go forward: %s", err.Error()))
			}
			page.WaitDOMStable(defaultWaitStableDur, defaultDomDiff)
			return mcp.NewToolResultText("Go forward successfully"), nil
		}
	}

	ReLoadHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to reload current page: %s", err.Error()))
			}
			err = page.Reload()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to reload current page: %s", err.Error()))
			}
			page.WaitDOMStable(defaultWaitStableDur, defaultDomDiff)
			return mcp.NewToolResultText("Reload current page successfully"), nil
		}
	}

	PressKeyHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to press key: %s", err.Error()))
			}
			keyStr := request.Params.Arguments["key"].(string)

			// 使用InsertText方法
			if len(keyStr) == 1 {
				err = page.InsertText(keyStr)
			} else {
				// 特殊键处理
				switch keyStr {
				case "Enter":
					err = page.Keyboard.Press(input.Enter)
				case "Tab":
					err = page.Keyboard.Press(input.Tab)
				case "Backspace":
					err = page.Keyboard.Press(input.Backspace)
				case "Escape":
					err = page.Keyboard.Press(input.Escape)
				case "ArrowUp":
					err = page.Keyboard.Press(input.ArrowUp)
				case "ArrowDown":
					err = page.Keyboard.Press(input.ArrowDown)
				case "ArrowLeft":
					err = page.Keyboard.Press(input.ArrowLeft)
				case "ArrowRight":
					err = page.Keyboard.Press(input.ArrowRight)
				default:
					// 尝试使用InsertText
					err = page.InsertText(keyStr)
				}
			}

			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to press key %s: %s", keyStr, err.Error()))
			}
			return mcp.NewToolResultText(fmt.Sprintf("Press key %s successfully", keyStr)), nil
		}
	}

	ClickHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to click element: %s", err.Error()))
			}
			selector := request.Params.Arguments["selector"].(string)
			element, err := page.Element(selector)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to find element %s: %s", selector, err.Error()))
			}
			err = element.Click(proto.InputMouseButtonLeft, 1)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to click element %s: %s", selector, err.Error()))
			}
			return mcp.NewToolResultText(fmt.Sprintf("Click element %s successfully", selector)), nil
		}
	}

	FillHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to fill out element: %s", err.Error()))
			}
			selector := request.Params.Arguments["selector"].(string)
			value := request.Params.Arguments["value"].(string)
			element, err := page.Element(selector)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to find element %s: %s", selector, err.Error()))
			}
			err = element.Input(value)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to fill out element %s: %s", selector, err.Error()))
			}
			return mcp.NewToolResultText(fmt.Sprintf("Fill out element %s successfully", selector)), nil
		}
	}

	DragHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to perform drag and drop: %s", err.Error()))
			}
			sourceSelector := request.Params.Arguments["source_selector"].(string)
			targetSelector := request.Params.Arguments["target_selector"].(string)

			source, err := page.Element(sourceSelector)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to find source element %s: %s", sourceSelector, err.Error()))
			}

			target, err := page.Element(targetSelector)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to find target element %s: %s", targetSelector, err.Error()))
			}

			// 获取源和目标元素的信息
			sourceShape, err := source.Shape()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to get shape of source element %s: %s", sourceSelector, err.Error()))
			}

			targetShape, err := target.Shape()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to get shape of target element %s: %s", targetSelector, err.Error()))
			}

			sourceBox := sourceShape.Box()
			targetBox := targetShape.Box()

			// 计算中心点
			sourcePoint := proto.NewPoint(sourceBox.X+sourceBox.Width/2, sourceBox.Y+sourceBox.Height/2)
			targetPoint := proto.NewPoint(targetBox.X+targetBox.Width/2, targetBox.Y+targetBox.Height/2)

			// 执行拖放操作
			mouse := page.Mouse

			// 移动到源元素中心
			err = mouse.MoveTo(sourcePoint)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to move mouse to source element: %s", err.Error()))
			}

			// 按下鼠标
			err = mouse.Down(proto.InputMouseButtonLeft, 1)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to press mouse button: %s", err.Error()))
			}

			// 移动到目标元素中心
			err = mouse.MoveTo(targetPoint)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to move mouse to target element: %s", err.Error()))
			}

			// 释放鼠标
			err = mouse.Up(proto.InputMouseButtonLeft, 1)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to release mouse button: %s", err.Error()))
			}

			return mcp.NewToolResultText(fmt.Sprintf("Drag and drop from %s to %s successfully", sourceSelector, targetSelector)), nil
		}
	}

	SelectorHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to select option: %s", err.Error()))
			}
			selector := request.Params.Arguments["selector"].(string)
			value := request.Params.Arguments["value"].(string)

			element, err := page.Element(selector)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to find select element %s: %s", selector, err.Error()))
			}

			// 通过执行JS来选择选项
			script := fmt.Sprintf(`() => {
				const el = this;
				el.value = "%s";
				el.dispatchEvent(new Event('change'));
				return true;
			}`, value)

			_, err = element.Eval(script)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to select option %s in element %s: %s", value, selector, err.Error()))
			}

			return mcp.NewToolResultText(fmt.Sprintf("Selected option %s in element %s successfully", value, selector)), nil
		}
	}

	WaitHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			waitTime := request.Params.Arguments["time"].(float64)
			if waitTime > 10 {
				waitTime = 10 // 限制最大等待时间为10秒
			}

			time.Sleep(time.Duration(waitTime) * time.Second)

			return mcp.NewToolResultText(fmt.Sprintf("Waited for %.1f seconds", waitTime)), nil
		}
	}

	ScreenshotHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to take screenshot: %s", err.Error()))
			}

			name := request.Params.Arguments["name"].(string)

			if selector, ok := request.Params.Arguments["selector"].(string); ok && selector != "" {
				element, err := page.Element(selector)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Failed to find element %s: %s", selector, err.Error()))
				}

				// 获取元素截图
				_, err = element.Screenshot(proto.PageCaptureScreenshotFormatPng, 100)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Failed to take screenshot of element %s: %s", selector, err.Error()))
				}
			} else {
				var width, height int
				if w, ok := request.Params.Arguments["width"].(float64); ok {
					width = int(w)
				} else {
					width = 800 // 默认宽度
				}

				if h, ok := request.Params.Arguments["height"].(float64); ok {
					height = int(h)
				} else {
					height = 600 // 默认高度
				}

				err = page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
					Width:  width,
					Height: height,
				})
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Failed to set viewport: %s", err.Error()))
				}

				// 获取页面截图
				_, err = page.Screenshot(false, nil)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Failed to take screenshot: %s", err.Error()))
				}
			}

			// 这里可以保存截图到文件，但需要根据项目实际需求实现
			// utils.SaveScreenshot(name, screenshotBytes)

			return mcp.NewToolResultText(fmt.Sprintf("Screenshot %s taken successfully", name)), nil
		}
	}

	PdfHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to generate PDF: %s", err.Error()))
			}

			filePath := request.Params.Arguments["file_path"].(string)
			fileName := request.Params.Arguments["file_name"].(string)

			_, err = page.PDF(&proto.PagePrintToPDF{
				PrintBackground: true,
			})

			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to generate PDF: %s", err.Error()))
			}

			// 同样，这里需要根据项目实际需求实现保存PDF的功能
			// utils.SavePDF(filePath, fileName, pdfBytes)

			return mcp.NewToolResultText(fmt.Sprintf("PDF saved to %s/%s", filePath, fileName)), nil
		}
	}

	CloseBrowserHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// 获取browser实例
			browser := rodCtx.GetBrowser()
			if browser == nil {
				return nil, errors.New("browser not launched")
			}

			err := browser.Close()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to close browser: %s", err.Error()))
			}

			return mcp.NewToolResultText("Browser closed successfully"), nil
		}
	}

	SnapshotHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to capture accessibility snapshot: %s", err.Error()))
			}

			// 获取页面快照，这里调用Rod的JS执行能力来获取页面可访问性树
			snapshotVal, err := page.Eval(`() => {
				function captureElementInfo(element, depth = 0, maxDepth = 10) {
					if (!element || depth > maxDepth) return null;
					
					let info = {
						tag: element.tagName.toLowerCase(),
						id: element.id || "",
						class: element.className || "",
						text: element.innerText || "",
						attributes: {},
						children: []
					};
					
					// 获取元素属性
					for (let attr of element.attributes) {
						info.attributes[attr.name] = attr.value;
					}
					
					// 递归获取子元素
					for (let child of element.children) {
						let childInfo = captureElementInfo(child, depth + 1, maxDepth);
						if (childInfo) {
							info.children.push(childInfo);
						}
					}
					
					return info;
				}
				
				return captureElementInfo(document.documentElement);
			}`)

			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to capture accessibility snapshot: %s", err.Error()))
			}

			// 返回JSON结果
			return mcp.NewToolResultText(fmt.Sprintf("Snapshot captured successfully: %v", snapshotVal.Value)), nil
		}
	}

	EvaluateHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to evaluate script: %s", err.Error()))
			}

			script := request.Params.Arguments["script"].(string)

			result, err := page.Eval(script)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to evaluate script: %s", err.Error()))
			}

			return mcp.NewToolResultText(fmt.Sprintf("Script evaluated successfully, result: %v", result.Value)), nil
		}
	}

	SelectElementHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to access page: %s", err.Error()))
			}

			selector := request.Params.Arguments["selector"].(string)

			getAttributes := false
			if attr, ok := request.Params.Arguments["get_attributes"].(bool); ok {
				getAttributes = attr
			}

			getText := true
			if txt, ok := request.Params.Arguments["get_text"].(bool); ok {
				getText = txt
			}

			// 使用JavaScript查询元素并获取信息
			script := fmt.Sprintf(`() => {
				const elements = Array.from(document.querySelectorAll("%s"));
				return elements.map(el => {
					const result = {
						tagName: el.tagName.toLowerCase(),
						id: el.id || "",
						className: el.className || ""
					};
					
					if (%t) {
						result.text = el.textContent || "";
					}
					
					if (%t) {
						result.attributes = {};
						Array.from(el.attributes).forEach(attr => {
							result.attributes[attr.name] = attr.value;
						});
					}
					
					if (el.tagName === "INPUT" || el.tagName === "TEXTAREA" || el.tagName === "SELECT") {
						result.value = el.value || "";
					}
					
					// 计算元素的位置和大小
					const rect = el.getBoundingClientRect();
					result.rect = {
						x: rect.x,
						y: rect.y,
						width: rect.width,
						height: rect.height
					};
					
					return result;
				});
			}`, selector, getText, getAttributes)

			result, err := page.Eval(script)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to select elements with selector %s: %s", selector, err.Error()))
			}

			// 构建结果文本
			resultText := fmt.Sprintf("Found elements matching selector '%s':\n", selector)
			resultText += fmt.Sprintf("%v", result.Value)

			return mcp.NewToolResultText(resultText), nil
		}
	}

	GetTextHandler = func(rodCtx *types.Context) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			page, err := rodCtx.EnsurePage()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to access page: %s", err.Error()))
			}

			// 获取参数
			selector := ""
			if sel, ok := request.Params.Arguments["selector"].(string); ok {
				selector = sel
			}

			trim := true
			if t, ok := request.Params.Arguments["trim"].(bool); ok {
				trim = t
			}

			includeHidden := false
			if ih, ok := request.Params.Arguments["include_hidden"].(bool); ok {
				includeHidden = ih
			}

			// 构建JavaScript脚本
			var script string
			if selector == "" {
				// 获取整个页面的文本
				script = fmt.Sprintf(`() => {
					const getVisibleText = (node) => {
						if (node.nodeType === Node.TEXT_NODE) {
							const text = node.textContent;
							const parentStyle = window.getComputedStyle(node.parentElement);
							if (%t || (parentStyle.display !== 'none' && parentStyle.visibility !== 'hidden')) {
								return %t ? text.trim() : text;
							}
							return '';
						}
						
						if (node.nodeType !== Node.ELEMENT_NODE) {
							return '';
						}
						
						const style = window.getComputedStyle(node);
						if (!%t && (style.display === 'none' || style.visibility === 'hidden')) {
							return '';
						}
						
						let text = '';
						for (const child of node.childNodes) {
							text += getVisibleText(child);
						}
						return %t ? text.trim() : text;
					};
					
					return getVisibleText(document.body);
				}`, includeHidden, trim, includeHidden, trim)
			} else {
				// 获取特定元素的文本
				script = fmt.Sprintf(`() => {
					const elements = document.querySelectorAll("%s");
					if (elements.length === 0) {
						return "No elements found matching selector: %s";
					}
					
					return Array.from(elements).map(el => {
						if (!%t && (window.getComputedStyle(el).display === 'none' || 
							window.getComputedStyle(el).visibility === 'hidden')) {
							return "[Hidden element]";
						}
						
						const text = el.textContent || "";
						return %t ? text.trim() : text;
					}).join("\n\n");
				}`, selector, selector, includeHidden, trim)
			}

			// 执行脚本
			result, err := page.Eval(script)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Failed to get text: %s", err.Error()))
			}

			textContent := fmt.Sprintf("%v", result.Value)
			if selector == "" {
				return mcp.NewToolResultText(fmt.Sprintf("Text content of the page:\n%s", textContent)), nil
			} else {
				return mcp.NewToolResultText(fmt.Sprintf("Text content of elements matching '%s':\n%s", selector, textContent)), nil
			}
		}
	}
)

var (
	CommonTools = []mcp.Tool{
		Navigation,
		GoBack,
		GoForward,
		ReLoad,
		PressKey,
		Click,
		Fill,
		Selector,
		Drag,
		Screenshot,
		Pdf,
		Wait,
		Snapshot,
		CloseBrowser,
		Evaluate,
		SelectElement,
		GetText,
	}
	CommonToolHandlers = map[string]ToolHandler{
		"rod_navigate":       NavigationHandler,
		"rod_go_back":        GoBackHandler,
		"rod_go_forward":     GoForwardHandler,
		"rod_reload":         ReLoadHandler,
		"rod_press_key":      PressKeyHandler,
		"rod_click":          ClickHandler,
		"rod_fill":           FillHandler,
		"rod_selector":       SelectorHandler,
		"rod_drag":           DragHandler,
		"rod_screenshot":     ScreenshotHandler,
		"rod_pdf":            PdfHandler,
		"rod_wait":           WaitHandler,
		"rod_snapshot":       SnapshotHandler,
		"rod_close":          CloseBrowserHandler,
		"rod_evaluate":       EvaluateHandler,
		"rod_select_element": SelectElementHandler,
		"rod_get_text":       GetTextHandler,
	}
)
