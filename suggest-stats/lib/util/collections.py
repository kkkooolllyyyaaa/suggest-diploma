class EffectiveList:
    def __init__(self, capacity):
        self.list = [[]] * capacity
        self.index = 0

    def append(self, element):
        if self.index == len(self.list):
            self.list.append(element)
            self.index += 1
            return

        self.list[self.index] = element
        self.index += 1

    def get(self):
        if self.index < len(self.list):
            return self.list[:self.index]
        return self.list
