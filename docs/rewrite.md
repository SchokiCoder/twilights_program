I initially wanted to use Rust with iced because the goals seemed simple enough
despite iced's immature state and i wanted to try it out.
At the moment scaling seems to be not implemented at all, which is kind of
important.
Trying to get the current scaling factor from your Application object just
returns 1.0 (literally a magic number returned from the get function).  
  
Worse than that are bugs.  
I can alomst not touch the font settings at all without the text vanishing from
existence.
Also setting the text color doesn't change the text's color.  
  
Time to let Go.  
