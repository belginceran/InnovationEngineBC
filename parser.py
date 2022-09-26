#Given a markdown file, this should break out the headings, paragraphs, executable commands etc. 

class Parser:
  

    def __init__(self, markdownFilepath):
        self.markdownFilepath = markdownFilepath
        self.markdownElements = []
        self.codeBlockType = '```'
        self.headingType = '#'
        self.paragraphType = 'p'

       

        

    def parseMarkdown(self):
        self.markdownFile = open(self.markdownFilepath)
        special_characters = ["#", "`", "~" ,"-" ]

        char = self.markdownFile.read(1)

        while char:
            if char == '#':
                #print("Found a heading! " + char)
                self.processHeading(char)
            elif char == '`':
                if self.checkForCodeBlock(char):
                    self.processCodeSample(char)

            # No specific tag, creating paragraphs and continuing to loop to a special char
            else:
                self.processParagraph(char)
                char = self.markdownFile.read(1)

        
        self.markdownFile.close()

    # Iterates through line adding each character to heading string. Also collects the heading type.
    # Creates a markdown element storing the subtype and text value. Appends to all markdown elements.
    def processHeading(self, char):
        text = ""
        headingCount = 0
        while char != '\n':
            if char == '#':
                headingCount += 1
           
            text += char
            char = self.markdownFile.read(1)
        
        subtype = "h" + str(headingCount)
        self.createAndAppendElement(self.headingType, subtype, text)

    # Iterates through text until 3 back ticks ``` are found signifying the end of code block
    def processCodeSample(self, char):
        command = ""
        subtype = ""
        endOfCodeBlock = False
        self.markdownFile.read(2)
        # Reading through initial backtick line to see what the subtype is
        while char != '\n':
            char = self.markdownFile.read(1)
            subtype += char

        while not endOfCodeBlock:
            if (char == '`'):
                if self.checkForCodeBlock(char):
                    endOfCodeBlock = True
                    # Read the remaining bash ticks
                    
                    self.markdownFile.read(2)
                    
                    

            else:
                command += char
            # Read all 3 back ticks
            char = self.markdownFile.read(1)

        subtype = subtype.strip()
        command = command.strip()
        self.createAndAppendElement(self.codeBlockType, subtype, command)
     
    # Iterates until we find a heading or back-tick. If heading is found a paragraph element
    # is created with the existing text, and 
    def processParagraph(self, char):
        paragraph = ""
        while char != '#' and char != '':
            if char == '`':
                if self.checkForCodeBlock(char):
                    self.createAndAppendElement(self.paragraphType, 'paragraph', paragraph.strip())
                    self.processCodeSample(char)
                    char = self.markdownFile.read(1)

            paragraph += char
            char = self.markdownFile.read(1)
        
        if char == '#':
            self.createAndAppendElement(self.paragraphType, 'paragraph', paragraph.strip())
            self.processHeading(char)


    def createAndAppendElement(self, type, subtype, value):
        element = MarkdownElement(subtype, value)
        self.markdownElements.append((type, element))
        pass

    # Helper function ran when we hit a backtick. It checks the next two characters to see
    # if they are also backticks. If so, we enter a code block and return true
    def checkForCodeBlock(self, char):
        currentPosition = self.markdownFile.tell()
        if self.markdownFile.read(1) == '`' and self.markdownFile.read(1) == '`':
            self.markdownFile.seek(currentPosition)
            return True
        else:
            self.markdownFile.seek(currentPosition)
            return False
    
    
    # If we want to add process for bold and italicized text
    def processBoldText(self,char):
        pass

    # If want to process dashes for hidden titles etc. 
    def processDash(self, char):
        pass

class MarkdownElement:

    def __init__(self, subtype, value):
        self.subtype = subtype
        self.value = value

