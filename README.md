# Important Note

Almost none of the following is implemented yet.  This README was mainly a
roadmap of sorts.  Behaviours outlined here will almost certainly be changed in
some way or outright removed, especially non-standard 6502 behaviour.

## A simple 6502 emulator

The purpose of this project is to have a working 6502 in software.

Cycle accuracy is not a goal of this project.  Instead, the focus will be on
making sure the instructions operate properly and sequentially (ie, pipelining
as seen on real hardware is ignored).

Special registers will be implemented to allow communication with the host
system.  Most notably the host's keyboard and console.  Additional APIs may be
defined later (networking? lol).

Similar to mappers on the NES, the emulator will have the option to use memory
bank switching.  This will probably only be one or two banking schemes that
require little to no setup.

File system IO.  Direct and/or via open/close file and map to work RAM?

Configurable work RAM?  Enabled/Disabled, WorkRAM = open file buffer? filename
would need to be written somewhere.

Fire an NMI at 60hz?  This could allow synchronizing and timing things with the
host system a bit easier.  Should probably be a configuration option (either
command line, register, or both).

### Keyboard 'API' (unimplemented)

- Lock/read (a-la NES controllers)
- Interrupt/status register based system
- A simple input buffer (FIFO style)

#### Lock/Read

This method would require the a sizeable chunk of memory addresses assigned to
just the keyboard (at least a full page), but it would probably be the simplest
to program for on the ASM side.

#### Interrupts and/or Status Registers

This might not be it's own separate input method, but rather a modification of
the Lock/Read method.  A register would be made available to check for a new
keyboard state (button pressed/released) and would not update the keyboard
state until the next frame/NMI (clocked at a configurable rate? 60hz default?).

#### Input Buffer

An input buffer method would need to properly handle modifier keys.  Maybe
signal that there are modifiers to read by setting the 7th bit and making flags
for which modifiers available in another register.

Registers:
- [R] Input in buffer
- [R] Modifier for last read input from buffer
- [W] Drop buffer

### Memory mapping / Bank Switching

~~Probably only two schemes: Full 32k banks or one fixed bank at $C000 and a
swapple bank of 16k starting at $8000.~~

~~Simple swapping mechanism: write the bank number to a register (eg. no shifting
like MMC1).~~

NES mapper support is fully implemented.  Currently, only three mappers are
written: FullRW, MMC1, and NROM.

### Filesystem IO & Raw File Access (unimplemented)

Open file by writing full path to a page between $4000 and $5FFF.
Once opened, fill work RAM with file's contents.  If more than 8k, allow bank
swapping?
Write/flush file contents on register write.
Close/drop file contents/changes on register write.
Open/reload file from disk on register write.

Implement a file system API to list files and directories, create and delete
files, change working directory, etc.

Registers:
- IO command/API
  - Open file (WRAM)
  - Reload file (WRAM). Used when other commands use WRAM for data.
  - Write file (WRAM)
  - Close file. Do not flush. A write would be required before this command to
    save data.
  - chdir
  - lsdir
  - create/delete file
- Swap WRAM page (read/write more of file or avoid overwriting contents of WRAM
  with command output)
- Write-only filename bank (half a page?)
- All data returned from IO command will be written to WRAM (directory
  listings, command errors, etc).
- Command return code.
