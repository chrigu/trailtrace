export const useFileValidation = () => {
  const validateMp4File = (file: File): boolean => {
    if (!file) {
      console.error('No file provided')
      return false
    }

    if (file.type !== "video/mp4" && !file.name.toLowerCase().endsWith(".mp4")) {
      console.error('Selected file is not an MP4')
      return false
    }

    return true
  }

  return {
    validateMp4File
  }
}