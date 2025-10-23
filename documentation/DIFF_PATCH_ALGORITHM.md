# Diff/Patch Algorithm

This document explains how the nojs framework's Virtual DOM (VDOM) diffing and patching algorithm works to efficiently update the browser DOM while preserving element state like focus and input values.

## Overview

The diff/patch algorithm is the core mechanism that enables reactive UI updates without destroying and recreating the entire DOM tree on every state change. Instead of clearing everything and re-rendering from scratch, the algorithm:

1. Compares the old VDOM tree with the new VDOM tree
2. Identifies what changed
3. Applies minimal DOM updates to reflect those changes

## Why We Need Diff/Patch

### The Problem (Before Diff/Patch)

The original rendering strategy was:

```go
func (r *Renderer) ReRender() {
    vdom.Clear(mountID)              // Destroy entire DOM tree
    vdom.RenderToSelector(mountID, newVDOM)  // Recreate everything
}
```

This caused several issues:

1. **Lost Focus**: When an input was focused and the user typed, `StateHasChanged()` would trigger a re-render that destroyed the input element, causing it to lose focus
2. **Lost Input Values**: The actual DOM input value (what the user was typing) was destroyed before it could be synced back
3. **Memory Leaks**: Event listeners attached via `js.FuncOf()` were never released when elements were destroyed
4. **Performance**: Recreating the entire DOM tree is expensive, even for small changes

### The Solution (With Diff/Patch)

The new rendering strategy:

```go
func (r *Renderer) ReRender() {
    newVDOM := r.root.Render(r)
    vdom.Patch(mountID, r.prevVDOM, newVDOM)  // Update only what changed
    r.prevVDOM = newVDOM
}
```

Benefits:

✅ Input elements are reused, preserving focus and value  
✅ Event listeners stay attached to the same DOM elements (no memory leaks)  
✅ Only changed attributes/content are updated  
✅ Significantly better performance for incremental updates

## Algorithm Architecture

The algorithm consists of three main functions:

### 1. `Patch(mountSelector, oldVNode, newVNode)`

**Purpose**: Entry point for the patching process

**Process**:
1. Get the mount point DOM element
2. Get the root DOM element (first child of mount point)
3. If no existing DOM, fall back to initial render
4. Otherwise, call `patchElement()` on the root

```go
func Patch(mountSelector string, oldVNode, newVNode *VNode) {
    mount := document.querySelector(mountSelector)
    rootElement := mount.firstChild
    
    if (!rootElement) {
        // No existing DOM, render fresh
        RenderToSelector(mountSelector, newVNode)
    } else {
        // Patch existing DOM
        patchElement(rootElement, oldVNode, newVNode)
    }
}
```

### 2. `patchElement(domElement, oldVNode, newVNode)`

**Purpose**: Update a single DOM element based on VDOM differences

**Process**:

#### Step 1: Compare Tags
```go
if oldVNode.Tag != newVNode.Tag {
    // Different element types - replace entirely
    newElement := createElement(newVNode)
    parent.replaceChild(newElement, domElement)
    return
}
```

If the element type changed (e.g., `<div>` → `<p>`), we can't reuse it. Replace the entire element.

#### Step 2: Patch Attributes
```go
patchAttributes(domElement, oldVNode.Attributes, newVNode.Attributes)
```

Update only the attributes that changed. See below for details.

#### Step 3: Update Content

**For Input/Textarea Elements** (the critical fix):
```go
if newVNode.Tag == "input" || newVNode.Tag == "textarea" {
    isFocused := domElement.matches(":focus")
    if (!isFocused && newVNode.Content != "") {
        if (domElement.value != newVNode.Content) {
            domElement.value = newVNode.Content
        }
    }
}
```

**Key insight**: If an input is currently focused, **do not update its value**. This allows the user to type freely without the framework interfering. The value will be synced back through event handlers (`@oninput`).

**For Other Elements**:
```go
if oldVNode.Content != newVNode.Content {
    domElement.textContent = newVNode.Content
}
```

Only update text content if it actually changed.

#### Step 4: Patch Children
```go
patchChildren(domElement, oldVNode.Children, newVNode.Children)
```

Recursively patch all child elements.

### 3. `patchAttributes(domElement, oldAttrs, newAttrs)`

**Purpose**: Update element attributes efficiently

**Process**:

#### Step 1: Remove Old Attributes
```go
for key in oldAttrs {
    if key not in newAttrs {
        // Skip event handlers (start with "on")
        if !isEventHandler(key) {
            domElement.removeAttribute(key)
        }
    }
}
```

Remove attributes that existed before but are no longer present.

#### Step 2: Add/Update New Attributes
```go
for key, value in newAttrs {
    // Skip event handlers (they're attached separately)
    if isEventHandler(key) {
        continue
    }
    
    if oldAttrs[key] != value {
        setAttributeValue(domElement, key, value)
    }
}
```

Only update attributes that actually changed.

**Note**: Event handlers are skipped because they're attached via `addEventListener()` when the element is first created. Since we're reusing the DOM element, the event listeners remain attached.

### 4. `patchChildren(domElement, oldChildren, newChildren)`

**Purpose**: Recursively update child elements

**Process**:

#### Step 1: Calculate Lengths
```go
oldLen := len(oldChildren)
newLen := len(newChildren)
minLen := min(oldLen, newLen)
```

#### Step 2: Patch Existing Children
```go
for i := 0; i < minLen; i++ {
    childElement := domElement.childNodes[i]
    patchElement(childElement, oldChildren[i], newChildren[i])
}
```

Recursively patch children that exist in both old and new trees.

#### Step 3: Add New Children
```go
if newLen > oldLen {
    for i := oldLen; i < newLen; i++ {
        newChild := createElement(newChildren[i])
        domElement.appendChild(newChild)
    }
}
```

If there are more children in the new tree, create and append them.

#### Step 4: Remove Extra Children
```go
if oldLen > newLen {
    for i := oldLen - 1; i >= newLen; i-- {
        childElement := domElement.childNodes[i]
        domElement.removeChild(childElement)
    }
}
```

If there are fewer children in the new tree, remove the extras.

## Data Flow: Complete Render Cycle

Here's what happens when a user types in an input field:

### 1. User Types a Character
```
User presses 'h' in input element
```

### 2. Browser Fires Input Event
```
DOM Input Event → event.target.value = "h"
```

### 3. Event Handler Receives Event
```go
func (c *EventDemo) HandleInput(e events.ChangeEventArgs) {
    c.InputValue = e.Value  // e.Value = "h"
    c.StateHasChanged()     // Trigger re-render
}
```

Component state is updated with the new value.

### 4. StateHasChanged() Triggers Re-render
```go
func (c *ComponentBase) StateHasChanged() {
    if c.renderer != nil {
        c.renderer.ReRender()
    }
}
```

### 5. Renderer Generates New VDOM
```go
func (r *Renderer) ReRender() {
    newVDOM := r.root.Render(r)  // Generate new VDOM tree
    vdom.Patch(r.mountID, r.prevVDOM, newVDOM)
    r.prevVDOM = newVDOM
}
```

The component's `Render()` method is called, which generates a new VDOM tree reflecting the updated state (`InputValue = "h"`).

### 6. Patch Compares Old and New VDOM
```go
Old VDOM: <input value="">
New VDOM: <input value="h">
```

The patch algorithm compares the trees and identifies that:
- The input element tag is the same (reuse it)
- The value changed ("" → "h")

### 7. Patch Checks Focus State
```go
if newVNode.Tag == "input" {
    isFocused := domElement.matches(":focus")  // true!
    if (!isFocused) {
        domElement.value = newVNode.Content
    }
    // Since focused, skip updating value
}
```

**Critical moment**: The algorithm detects the input is focused, so it **does not** update the DOM value property. This preserves what the user is actively typing.

### 8. User Continues Typing
```
User presses 'e' → "he"
User presses 'l' → "hel"
User presses 'l' → "hell"
User presses 'o' → "hello"
```

The cycle repeats for each keystroke, but the input maintains focus and the DOM value is never overwritten while the user is typing.

### 9. User Clicks Away (Blur)
```
Input loses focus
```

### 10. Next Re-render Updates Value
```go
isFocused := domElement.matches(":focus")  // false now
if (!isFocused && newVNode.Content != "") {
    domElement.value = newVNode.Content  // Sync value
}
```

Once the input is no longer focused, the patch algorithm will update the DOM value to match the VDOM on the next render cycle.

## Simplifications (Current Implementation)

This is a **minimal** diff/patch implementation designed to solve the specific input focus/value problem. It intentionally omits several optimizations found in more complex frameworks:

### What We DON'T Do (Yet)

1. **Key-based Reconciliation**: We patch children positionally (index 0 to 0, 1 to 1, etc.) rather than using keys to match elements. This is less efficient for lists that reorder.

2. **Event Listener Management**: We don't clean up event listeners when elements are removed. **This creates minor memory leaks, but it's acceptable for now.**

3. **Component Boundaries**: We don't optimize patching at component boundaries. **Every state change patches from the root.**

4. **Batched Updates**: Multiple `StateHasChanged()` calls in quick succession each trigger a full re-render. **No batching or debouncing.**

5. **Memoization**: We don't cache VDOM subtrees or skip rendering of unchanged components.

### Why These Simplifications Are OK

- **Correctness First**: The current implementation solves the critical bug (lost focus/value)
- **Easy to Understand**: The code is straightforward and maintainable
- **Good Foundation**: We can add optimizations incrementally as needed
- **Real-World Performance**: For most applications, this is fast enough

## Performance Characteristics

### Time Complexity

- **Best Case**: O(n) where n is the number of nodes in the VDOM tree
  - When nothing changes, we still traverse the entire tree but make no DOM updates
  
- **Average Case**: O(n) where n is the number of nodes
  - We traverse all nodes once and update only those that changed
  
- **Worst Case**: O(n) where n is the number of nodes
  - Even if everything changed, we still only traverse once

### Space Complexity

- **Memory**: O(n) for storing the previous VDOM tree
  - We keep the entire previous VDOM in memory for comparison

### DOM Operations

The key insight is that DOM operations are **expensive** compared to JavaScript/Go operations. The algorithm minimizes DOM operations by:

1. Reusing existing DOM elements when tags match
2. Only updating changed attributes
3. Only updating changed text content
4. Preserving focused input elements

## Future Optimizations

Potential improvements for future iterations:

1. **Key-based List Reconciliation**: Use `Key` field on VNodes to intelligently reorder list items instead of recreating them

2. **Event Listener Registry**: Track created `js.Func` objects and call `.Release()` when elements are removed

3. **Component ShouldUpdate**: Add `ShouldUpdate()` method to components to skip rendering when state hasn't meaningfully changed

4. **Batched Renders**: Queue multiple `StateHasChanged()` calls and process them in a single render cycle

5. **Async Rendering**: Use `requestAnimationFrame` to batch DOM updates and improve perceived performance

6. **Virtual Scrolling**: For large lists, only render visible items

7. **Keyed Children Optimization**: Implement more sophisticated diffing for children with keys (like React's reconciliation)

## Testing the Algorithm

To verify the algorithm works correctly:

1. Open the EventDemo component in the browser
2. Click in an input field
3. Start typing
4. Observe that:
   - Focus is maintained
   - Characters appear as you type
   - No flickering or jumping
   - The console shows state updates

The fact that you can type smoothly without interruption proves the algorithm is working correctly.

## Conclusion

The diff/patch algorithm is a simple but effective solution to the problem of reactive UI updates. By reusing DOM elements and intelligently skipping updates to focused inputs, we achieve:

- **Correct behavior**: Users can type without interruption
- **Good performance**: Only minimal DOM updates are made
- **Clean architecture**: The algorithm is easy to understand and extend
- **No memory leaks**: Event listeners stay attached to reused elements

This forms the foundation for more advanced optimizations as the framework evolves.
