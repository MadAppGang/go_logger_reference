# go_logger_reference
The repo is a demonstratoin how to orginise loggers across the system

### StreamProcessor
Service processing stread of data.

Focused on separate context for each processing stage rather than for each processed data unit.

Logger can be created inside each processing component.

### RequestProcessor
Service focused mainly on request processing.

Context of processing each separate data unit(request) is more important than context of processing component.
Processing of each data unit have to be tracked throughout all the stages.

Logger probably should be created at the point where data unit processing starts and passed downstack along with unit being processed.
